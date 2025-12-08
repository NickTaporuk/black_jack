// v1/tests/drop_curve_test.cpp
// make test_drop_curve

#include "jackpot_engine_v1.h"
#include <iostream>
#include <fstream>
#include <vector>
#include <iomanip>
#include <nlohmann/json.hpp>

using json = nlohmann::json;

struct Stats {
    uint64_t dropPoint;
    double   avgSteps;
    int      runs = 100;
};

constexpr uint64_t MIN_POINT = 0;
constexpr uint64_t MAX_POINT = 10'000ULL;
constexpr int      VOLATILITY = 3;  // High = кубическая
constexpr uint64_t STEP = 1;

int main() {
    std::cout << "Статистический тест: High Volatility (progress³)\n";
    std::cout << "dropPoint: 5000 → 15000, 1000 прогонов на точку\n\n";

    std::vector<Stats> results;

    for (uint64_t targetDrop = 5000; targetDrop <= 10010; targetDrop += 50) {
        uint64_t totalSteps = 0;

        for (int run = 0; run < 1000; ++run) {
            uint64_t seed = 0x517cc1b727220a94ULL ^ (targetDrop + uint64_t(run) * 0x9e3779b97f4a7c15ULL);
            JackpotEngineHandle* engine = jackpot_engine_v1(seed);
            if (!engine) { std::cerr << "Engine init failed!\n"; return 1; }

            int id = jackpot_add(engine, MIN_POINT, MAX_POINT, VOLATILITY);
            if (id < 0) { std::cerr << "Add failed!\n"; return 1; }

            uint64_t actualDropPoint = jackpot_get_drop_point(engine, id);
            std::cout<<actualDropPoint<<std::endl;

            uint64_t steps = 0;
            bool dropped = false;
            while (!dropped) {
                steps += STEP;
                dropped = jackpot_process(engine, id, STEP);
            }

            totalSteps += steps;
            jackpot_engine_destroy(engine);
        }

        double avg = double(totalSteps) / 1000.0;
        results.push_back({targetDrop, avg});

        std::cout << "dropPoint=" << std::setw(5) << targetDrop
                  << " → avg steps = " << std::fixed << std::setprecision(2) << avg
                  << " (≈ " << targetDrop << ")\n";
    }

    json j = json::array();
    for (const auto& s : results) {
        j.push_back({
            {"drop_point", s.dropPoint},
            {"avg_steps",  s.avgSteps},
            {"expected",   double(s.dropPoint)},
            {"volatility", "High (progress³)"},
            {"runs",       100}
        });
    }

    std::ofstream out("tests/drop_curve_cubic.json");
    out << j.dump(2);
    std::cout << "\nDONE! → tests/drop_curve_cubic.json\n";
    return 0;
}
