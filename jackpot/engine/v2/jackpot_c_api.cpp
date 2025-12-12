#include "jackpot_c_api.h"
#include "jackpot_engine.hpp"

int jackpot_check_uniform(
    uint64_t seed,
    uint64_t min_point,
    uint64_t max_point,
    uint64_t current,
    uint64_t* out_drop_point
) {
    if (!out_drop_point || min_point >= max_point) {
        return 0;
    }

    JackpotEngine eng(seed);
    uint64_t dp = eng.drop(min_point, max_point);
    *out_drop_point = dp;

    return JackpotEngine::isDropped(current, dp) ? 1 : 0;
}
