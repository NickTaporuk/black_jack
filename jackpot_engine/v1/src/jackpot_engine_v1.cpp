#include "jackpot_engine_v1.h"
#include "jackpot_engine_core.hpp"

extern "C" {

struct JackpotEngineHandle {
    JackpotEngine* engine;
};

JackpotEngineHandle* jackpot_engine_v1(uint64_t seed) {
    try {
        auto rng = std::make_unique<MtRng>(seed);
        JackpotEngine* e = new JackpotEngine(std::move(rng));
        return new JackpotEngineHandle{e};
    } catch (...) {
        return nullptr;
    }
}

int jackpot_add(JackpotEngineHandle* h,
                uint64_t minPoint,
                uint64_t maxPoint,
                int volatility)
{
    if (!h || !h->engine) return -1;
    try {
        JackpotConfig cfg{minPoint, maxPoint, static_cast<Volatility>(volatility)};
        return h->engine->addJackpot(cfg);
    } catch (...) {
        return -1;
    }
}

int jackpot_process(JackpotEngineHandle* h, int id, uint64_t step) {
    if (!h || !h->engine) return -1;
    try {
        return h->engine->processJackpot(id, step) ? 1 : 0;
    } catch (...) {
        return -1;
    }
}

void jackpot_reset(JackpotEngineHandle* h, int id) {
    if (h && h->engine) {
        try { h->engine->resetJackpot(id); } catch (...) {}
    }
}

void jackpot_enable(JackpotEngineHandle* h, int id) {
    if (h && h->engine) {
        try { h->engine->enableJackpot(id); } catch (...) {}
    }
}

void jackpot_disable(JackpotEngineHandle* h, int id) {
    if (h && h->engine) {
        try { h->engine->disableJackpot(id); } catch (...) {}
    }
}

uint64_t jackpot_get_counter(JackpotEngineHandle* h, int id) {
    if (!h || !h->engine) return 0;
    try { return h->engine->getCounter(id); } catch (...) { return 0; }
}

uint64_t jackpot_get_drop_point(JackpotEngineHandle* h, int id) {
    if (!h || !h->engine) return 0;
    try { return h->engine->getDropPoint(id); } catch (...) { return 0; }
}

int jackpot_is_active(JackpotEngineHandle* h, int id) {
    if (!h || !h->engine) return 0;
    try { return h->engine->isActive(id) ? 1 : 0; } catch (...) { return 0; }
}

size_t jackpot_count(JackpotEngineHandle* h) {
    if (!h || !h->engine) return 0;
    try { return h->engine->count(); } catch (...) { return 0; }
}

void jackpot_engine_destroy(JackpotEngineHandle* h) {
    if (!h) return;
    delete h->engine;
    delete h;
}

} // extern "C"
