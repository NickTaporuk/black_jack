#include "../src/jackpot_engine_v1.h"
#include <iostream>
#include <fstream>
#include <vector>
#include <chrono>
#include <nlohmann/json.hpp>

using json = nlohmann::json;
using namespace std::chrono;

struct HitRecord {
    int id;
    uint64_t counter;
    uint64_t dropPoint;
    double progress() const { return double(counter) / dropPoint; }
};

int main() {
    JackpotEngineHandle* eng = jackpot_engine_v1(123456789ULL);
    if (!eng) return 1;

    struct Cfg { uint64_t min, max; int v; };
    Cfg cfgs[] = {
        {500, 1500, 1},
        {10000, 50000, 2},
        {500000, 2000000, 3},
    };

    int ids[3];
    for (int i = 0; i < 3; i++) {
        ids[i] = jackpot_add(eng, cfgs[i].min, cfgs[i].max, cfgs[i].v);
    }

    const uint64_t N = 1ULL;
    std::vector<HitRecord> hits;

    for (uint64_t step = 1; step <= N; step++) {
        for (int i = 0; i < 3; i++) {
            if (jackpot_process(eng, ids[i], 1) == 1) {
                hits.push_back({
                    ids[i],
                    jackpot_get_counter(eng, ids[i]),
                    jackpot_get_drop_point(eng, ids[i])
                });
            }
        }
    }

    json rep;
    rep["hits"] = json::array();
    for (auto& h : hits) {
        rep["hits"].push_back({
            {"id", h.id},
            {"counter", h.counter},
            {"drop_point", h.dropPoint},
            {"progress", h.progress()}
        });
    }

    std::ofstream("v1/tests/simulation_report.json") << rep.dump(2);
    jackpot_engine_destroy(eng);

    std::cout << "Simulation done!" << std::endl;
}
