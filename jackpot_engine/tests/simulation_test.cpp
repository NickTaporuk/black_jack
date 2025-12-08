// v1/tests/simulation_test.cpp
// Build: make test_simulation
// Output: simulation_report.json + performance.txt

#include "../src/jackpot_engine_v1.h"
#include <iostream>
#include <fstream>
#include <vector>
#include <chrono>
#include <nlohmann/json.hpp>

using json = nlohmann::json;
using namespace std::chrono;

struct HitRecord {
    int jackpot_id;
    uint64_t counter_at_hit;
    uint64_t drop_point;
    double progress() const { return static_cast<double>(counter_at_hit) / drop_point; }
};

int main() {
    std::cout << "Starting certification-grade jackpot engine simulation...\n";

    //auto start_total = high_resolution_clock::now();

    JackpotEngineHandle* engine = jackpot_engine_v1(0x517cc1b727220a94ULL);
    if (!engine) {
        std::cerr << "Failed to create engine!\n";
        return 1;
    }

    struct Config { uint64_t min, max; int vol; };
    Config configs[] = {
        {500,      1500,     1},   // Mini  – Low volatility
        {10000,    50000,    2},   // Major – Medium volatility
        {500000,   2000000,  3}    // Mega  – High volatility
    };

    int ids[3];
    for (int i = 0; i < 3; ++i) {
        ids[i] = jackpot_add(engine, configs[i].min, configs[i].max, configs[i].vol);
        if (ids[i] < 0) {
            std::cerr << "Failed to create jackpot!\n";
            return 1;
        }
    }

    const uint64_t TOTAL_BETS = 100'000'000ULL;
    std::vector<HitRecord> hits;

    auto start_perf = high_resolution_clock::now();

    for (uint64_t bet = 1; bet <= TOTAL_BETS; ++bet) {
        for (int i = 0; i < 3; ++i) {
            int dropped = jackpot_process(engine, ids[i], 1);
            if (dropped) {
                uint64_t counter = jackpot_get_counter(engine, ids[i]);
                uint64_t drop_pt = jackpot_get_drop_point(engine, ids[i]);
                hits.push_back({ids[i], counter, drop_pt});
            }
        }
    }

    auto end_perf = high_resolution_clock::now();
    double seconds = duration<double>(end_perf - start_perf).count();
    double rps = TOTAL_BETS / seconds;

    json report;
    report["simulation"] = {
        {"total_bets", TOTAL_BETS},
        {"duration_seconds", seconds},
        {"requests_per_second", rps},
        {"hits_count", hits.size()}
    };

    json jackpot_stats = json::array();
    for (int i = 0; i < 3; ++i) {
        jackpot_stats.push_back({
            {"id", ids[i]},
            {"min_point", configs[i].min},
            {"max_point", configs[i].max},
            {"volatility", configs[i].vol}
        });
    }
    report["jackpots"] = jackpot_stats;

    json hits_json = json::array();
    for (const auto& h : hits) {
        hits_json.push_back({
            {"jackpot_id", h.jackpot_id},
            {"counter_at_hit", h.counter_at_hit},
            {"drop_point", h.drop_point},
            {"progress", h.progress()}
        });
    }
    report["hits"] = hits_json;

    std::ofstream out("v1/tests/simulation_report.json");
    out << report.dump(2);

    std::ofstream perf("v1/tests/performance.txt");
    perf << "Certification simulation completed successfully\n";
    perf << "Total bets processed: " << TOTAL_BETS << "\n";
    perf << "Execution time: " << seconds << " seconds\n";
    perf << "Processing speed: " << static_cast<uint64_t>(rps) << " bets/second\n";
    perf << "Total jackpot hits: " << hits.size() << "\n";

    std::cout << "Done! performance.txt and simulation_report.json created.\n";

    jackpot_engine_destroy(engine);
    return 0;
}