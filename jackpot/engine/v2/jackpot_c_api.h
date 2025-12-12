#pragma once
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

// Возвращает 1 если дропнулся, 0 если нет
int jackpot_check_uniform(
    uint64_t seed,
    uint64_t min_point,
    uint64_t max_point,
    uint64_t current,
    uint64_t* out_drop_point
);

#ifdef __cplusplus
}
#endif
