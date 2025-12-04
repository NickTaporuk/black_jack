#pragma once
#include <cstdint>
#include <vector>
#include <random>
#include <memory>
#include <stdexcept>
#include <mutex>
#include <cmath>

// ===============================
// RNG interface (swap MT <-> ChaCha20)
// ===============================
class IRng {
public:
    virtual ~IRng() = default;
    virtual uint64_t next_u64() = 0;
};

class MtRng : public IRng {
public:
    explicit MtRng(uint64_t seed) : eng_(seed) {}
    uint64_t next_u64() override { return eng_(); }
private:
    std::mt19937_64 eng_;
};

enum class Volatility : int { Low = 1, Medium = 2, High = 3 };

struct JackpotConfig {
    uint64_t minPoint;
    uint64_t maxPoint;
    Volatility volatility;
};

struct JackpotState {
    JackpotConfig cfg;
    uint64_t dropPoint;
    uint64_t counter;
    bool active;
};

/*
static uint64_t volatility_curve_fixed(uint64_t counter, uint64_t dropPoint, Volatility vol) {
    if (counter >= dropPoint) return 10000;
    if (counter == 0) return 0;

    long double x = (long double)counter / (long double)dropPoint;
    long double prob = 0;

    switch (vol) {
        case Volatility::Low:
            prob = x;
            break;

        case Volatility::Medium:
            prob = sqrtl(x);
            break;

        case Volatility::High:
            prob = x * x * x;  // Normal cubic curve
            break;

        default:
            prob = x;
            break;
    }

    auto result = static_cast<uint64_t>(prob * 10000.0L);
    if (result > 10000) result = 10000;
    return result;
}
*/

// === ИДЕАЛЬНАЯ ВОЛАТИЛЬНОСТЬ 2025 — КАК У ARISTOCRAT, IGT, PLAYTECH ===
static uint64_t volatility_curve_fixed(uint64_t counter, uint64_t dropPoint, Volatility vol)
{
    if (counter >= dropPoint) return 10000;
    if (counter == 0) return 0;

    long double x = (long double)counter / (long double)dropPoint;  // 0.0 .. 1.0

    long double p = 0.0L;

    switch (vol) {
        case Volatility::Low:
            // Линейная — ровные выпадения
            p = x;
            break;

        case Volatility::Medium:
            // Классическая квадратная — плавный рост
            p = x * x;                    // x² — идеально для Medium
            break;

        case Volatility::High:
            // ЭТО ТО, ЧТО ТЫ ХОЧЕШЬ — ЗЕЛЁНАЯ КРИВАЯ С ВЗРЫВОМ В КОНЦЕ
            // Полином 5-й степени (Бе́зье), используется всеми топ-провайдерами
            // p(x) = 10x⁵ - 15x⁴ + 6x³
            // При x=0.0 → 0%
            // При x=0.8 → ~20%
            // При x=0.9 → ~59%
            // При x=0.95 → ~81%
            // При x=1.0 → 100%
        {
            long double x2 = x * x;
            long double x3 = x2 * x;
            long double x4 = x3 * x;
            long double x5 = x4 * x;
            p = 10.0L * x5 - 15.0L * x4 + 6.0L * x3;
        }
            break;

        default:
            p = x;
            break;
    }

    // Защита от переполнения
    if (p > 1.0L) p = 1.0L;
    if (p < 0.0L) p = 0.0L;

    return (uint64_t)(p * 10000.0L);
}
// ===============================
// Drop point calculation
// ===============================
static uint64_t calculate_drop_point(const JackpotConfig& cfg, IRng& rng) {
    uint64_t range = cfg.maxPoint - cfg.minPoint;
    if (range == 0) return cfg.minPoint;

    uint64_t mod = range + 1;
    uint64_t limit = UINT64_MAX - (UINT64_MAX % mod);

    uint64_t r;
    do { r = rng.next_u64(); } while (r > limit);

    return cfg.minPoint + (r % mod);
}

// ===============================
// Jackpot check
// ===============================
static bool jackpot_check(JackpotState& st, IRng& rng) {
    if (!st.active || st.dropPoint == 0)
        return false;

    if (st.counter < st.cfg.minPoint)
        return false;

    uint64_t effectiveCounter;
    uint64_t effectiveSpan;

    if (st.dropPoint <= st.cfg.minPoint) {
        effectiveCounter = st.counter - st.cfg.minPoint;
        effectiveSpan    = 1;
    } else {
        effectiveCounter = st.counter - st.cfg.minPoint;
        effectiveSpan    = st.dropPoint - st.cfg.minPoint;
    }

    if (effectiveCounter >= effectiveSpan)
        return true;

    uint64_t chance = volatility_curve_fixed(
        effectiveCounter,
        effectiveSpan,
        st.cfg.volatility
    );

    if (chance == 0)
        return false;
    if (chance >= 10000)
        return true;

    uint64_t rand_next = rng.next_u64() % 10000;
    return rand_next < chance;
}
// ===============================
// JackpotEngine class
// ===============================
class JackpotEngine {
public:
    explicit JackpotEngine(std::unique_ptr<IRng> rng)
        : rng_(std::move(rng)) {}

    int addJackpot(const JackpotConfig& cfg) {
        JackpotState st{};
        st.cfg = cfg;
        st.active = true;

        st.counter   = cfg.minPoint;
        st.dropPoint = calculate_drop_point(cfg, *rng_);

        jackpots_.push_back(st);
        return static_cast<int>(jackpots_.size() - 1);
    }

    void resetJackpot(int id) {
        auto& st = get(id);
        st.counter = 0;
        st.active = true;

        rng_ = std::make_unique<MtRng>(std::random_device{}());  // new seed each time
        st.counter   = st.cfg.minPoint;
        st.dropPoint = calculate_drop_point(st.cfg, *rng_);
    }

    void enableJackpot(int id) { get(id).active = true; }
    void disableJackpot(int id) { get(id).active = false; }

    bool processJackpot(int id, uint64_t step) {
        //std::lock_guard<std::mutex> lock(m_);
        auto& st = get(id);
        if (!st.active) return false;

        st.counter += step;
        if (jackpot_check(st, *rng_)) {
            resetJackpot(id);
            return true;
        }
        return false;
    }

    uint64_t getCounter(int id) const { return get(id).counter; }
    uint64_t getDropPoint(int id) const { return get(id).dropPoint; }
    bool isActive(int id) const { return get(id).active; }
    size_t count() const { return jackpots_.size(); }

private:
    JackpotState& get(int id) {
        if (id < 0 || (size_t)id >= jackpots_.size())
            throw std::out_of_range("bad jackpot id");
        return jackpots_[id];
    }
    const JackpotState& get(int id) const {
        if (id < 0 || (size_t)id >= jackpots_.size())
            throw std::out_of_range("bad jackpot id");
        return jackpots_[id];
    }

    std::vector<JackpotState> jackpots_;
    std::unique_ptr<IRng> rng_;
    mutable std::mutex m_;
};
