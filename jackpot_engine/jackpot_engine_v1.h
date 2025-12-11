#pragma once
#include <stdint.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct JackpotEngineHandle JackpotEngineHandle;

JackpotEngineHandle* jackpot_engine_v1(uint64_t seed);

int jackpot_add(JackpotEngineHandle* h,
                uint64_t minPoint,
                uint64_t maxPoint,
                int volatility);

int jackpot_process(JackpotEngineHandle* h,
                    int jackpot_id,
                    uint64_t step_amount);

void jackpot_reset(JackpotEngineHandle* h, int jackpot_id);

void jackpot_enable(JackpotEngineHandle* h, int jackpot_id);
void jackpot_disable(JackpotEngineHandle* h, int jackpot_id);

uint64_t jackpot_get_counter(JackpotEngineHandle* h, int jackpot_id);
uint64_t jackpot_get_drop_point(JackpotEngineHandle* h, int jackpot_id);
int      jackpot_is_active(JackpotEngineHandle* h, int jackpot_id);
size_t   jackpot_count(JackpotEngineHandle* h);

void jackpot_engine_destroy(JackpotEngineHandle* h);

const char* jackpot_engine_version();

#ifdef __cplusplus
}
#endif
