#pragma once
#include <cstdint>
#include <vector>
#include <random>
#include <memory>
#include <stdexcept>
#include <mutex>
#include <cmath>
#include <unistd.h>
#include <ostream>

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

inline std::ostream& operator<<(std::ostream& os, Volatility v) {
    switch (v) {
        case Volatility::Low:    os << "Low"; break;
        case Volatility::Medium: os << "Medium"; break;
        case Volatility::High:   os << "High"; break;
        default:                 os << "Unknown"; break;
    }
    return os;
}

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

/*static uint64_t volatility_curve_fixed(uint64_t counter, uint64_t span, Volatility vol) {
    if (span == 0) {
        return 10000;
    }

    if (counter >= span) {
        return 10000;
    }

    uint64_t k = counter + 1;
    uint64_t N = span;

    long double alpha;
    switch (vol) {
        case Volatility::Low:    alpha = 1.5L; break;  // мягкая кривая
        case Volatility::Medium: alpha = 3.0L; break;  // стандартно
        case Volatility::High:   alpha = 6.0L; break;  // сильно в конец
        default:                 alpha = 3.0L; break;
    }

    long double Nd   = static_cast<long double>(N);
    long double x_k   = static_cast<long double>(k)       / Nd;
    long double x_k_1 = static_cast<long double>(k - 1)   / Nd;

    long double F_k   = powl(x_k,   alpha);
    long double F_k_1 = (k == 1) ? 0.0L : powl(x_k_1, alpha);

    long double numerator   = F_k - F_k_1;
    long double denominator = 1.0L - F_k_1;

    long double q_k;
    if (denominator <= 0.0L) {
        // В конце отрезка F ≈ 1 → вынужденно q_k = 1
        q_k = 1.0L;
    } else {
        q_k = numerator / denominator;
    }

    if (q_k <= 0.0L) return 0;
    if (q_k >= 1.0L) return 10000;

    uint64_t result = static_cast<uint64_t>(q_k * 10000.0L);
    if (result > 10000) result = 10000;
    return result;
}*/

static uint64_t volatility_curve_fixed(uint64_t counter,
                                       uint64_t span,
                                       Volatility vol)
{
    if (span == 0) return 10000;
    if (counter >= span) return 10000;

    // k in [0 .. span-1]
    long double x = (long double)counter / (long double)(span - 1);

    // position of the peak
    long double peak;
    switch (vol) {
        case Volatility::Low:    peak = 0.40L; break;
        case Volatility::Medium: peak = 0.50L; break;
        case Volatility::High:   peak = 0.65L; break;
        default:                 peak = 0.50L; break;
    }

    long double w;
    if (x <= peak) {
        w = x / peak;
    } else {
        w = (1.0L - x) / (1.0L - peak);
    }

    if (w <= 0.0L) return 0;
    if (w >= 1.0L) return 10000;

    return static_cast<uint64_t>(w * 10000.0L);
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

static bool jackpot_check(JackpotState& st, IRng& rng) {
    if (!st.active)
        return false;

    if (st.counter < st.cfg.minPoint)
        return false;

    uint64_t span = st.cfg.maxPoint - st.cfg.minPoint + 1;
    uint64_t k    = st.counter - st.cfg.minPoint + 1;
    if (k > span) k = span;

    long double alpha;
    switch (st.cfg.volatility) {
        case Volatility::Low:    alpha = 1.5L; break;
        case Volatility::Medium: alpha = 3.0L; break;
        case Volatility::High:   alpha = 6.0L; break;
        default:                 alpha = 3.0L; break;
    }

    long double N     = (long double)span;
    long double x_k   = (long double)k       / N;
    long double x_k_1 = (long double)(k - 1) / N;

    long double F_k   = powl(x_k,   alpha);
    long double F_k_1 = (k == 1) ? 0.0L : powl(x_k_1, alpha);

    long double numerator   = F_k - F_k_1;
    long double denominator = 1.0L - F_k_1;

    long double q_k;
    if (denominator <= 0.0L) {
        q_k = 1.0L;
    } else {
        q_k = numerator / denominator;
    }

    if (q_k <= 0.0L)
        return false;
    if (q_k >= 1.0L)
        return true;

    uint64_t u = rng.next_u64();
    long double u01 = (long double)u / (long double)UINT64_MAX;
    return (u01 < q_k);
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

static inline uint64_t mix_u64(uint64_t x) {
    // splitmix64 — прекрасный быстрый хеш для энтропии
    x += 0x9e3779b97f4a7c15ULL;
    x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9ULL;
    x = (x ^ (x >> 27)) * 0x94d049bb133111ebULL;
    return x ^ (x >> 31);
}

void resetJackpot(int id) {
    auto& st = get(id);

    st.counter = 0;
    st.active = true;

    // --- FORTIFIED ENTROPY SOURCE ---
    uint64_t t = std::chrono::high_resolution_clock::now().time_since_epoch().count();
    uint64_t pid = (uint64_t)getpid();

    uint64_t rd1 = 0, rd2 = 0;
    try {
        std::random_device rd;
        rd1 = ((uint64_t)rd() << 32) ^ rd();
        rd2 = ((uint64_t)rd() << 32) ^ rd();
    } catch (...) {
        rd1 = 0x123456789ABCDEFULL;
        rd2 = 0xCAFEBABEDEADBEEFULL;
    }

    uint64_t seed = mix_u64(t) ^ mix_u64(pid) ^ mix_u64(rd1) ^ mix_u64(rd2);

    // install new RNG with proper seed
    rng_ = std::make_unique<MtRng>(seed);

    // initialize fields
    st.counter = st.cfg.minPoint;
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
