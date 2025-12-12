#pragma once
#include <cstdint>
#include <random>
#include <stdexcept>

class JackpotEngine {
public:
    explicit JackpotEngine(uint64_t seed) : rng_(seed) {}

    uint64_t drop(uint64_t minPoint, uint64_t maxPoint) {
        if (minPoint >= maxPoint) {
            throw std::invalid_argument("minPoint must be < maxPoint");
        }
        return uniform_range(minPoint, maxPoint);
    }

    static bool isDropped(uint64_t current, uint64_t dropPoint) {
        return current >= dropPoint;
    }

private:
    uint64_t uniform_range(uint64_t min, uint64_t max) {
        uint64_t span  = max - min + 1;
        uint64_t limit = UINT64_MAX - (UINT64_MAX % span);
        uint64_t r;
        do { r = rng_(); } while (r > limit);
        return min + (r % span);
    }

    std::mt19937_64 rng_;
};
