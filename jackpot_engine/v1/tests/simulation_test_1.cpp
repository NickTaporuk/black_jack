#include <iostream>
#include <thread>
#include "../src/jackpot_engine_v1.h"
#include "../src/jackpot_engine_core.hpp"
#include <chrono>
#include "../src/jackpot_engine_v1.cpp"

int main() {
    std::cout << "Run simulation...\n";

    auto* eng = jackpot_engine_v1(std::chrono::high_resolution_clock::now().time_since_epoch().count());
    if (!eng) {
        std::cerr << "Engine init failed\n";
        return 1;
    }

    int id = jackpot_add(eng, 19999980, 20000000, 3);
    if (id < 0) {
        std::cerr << "Jackpot add failed\n";
        return 1;
    }

    std::cout << "Jackpot created with id=" << id << "\n";

    uint64_t dp = jackpot_get_drop_point(eng, id);
    std::cout << "dropPoint=" << dp << "\n";

    uint64_t step = 1;
    uint64_t iterations = 0;

    while (true) {
        int hit = jackpot_process(eng, id, step);
        iterations++;

        if (iterations % 100 == 0) {
            std::cout << "Progress: iterations=" << iterations
                      << " counter=" << jackpot_get_counter(eng, id)
                      << " / dropPoint=" << jackpot_get_drop_point(eng, id)
                      << "\n";
        }

        if (hit == 1) {
            std::cout << "\n🎉 JACKPOT HIT after " << iterations << " iterations\n";
            std::cout << "New dropPoint=" << jackpot_get_drop_point(eng, id) << "\n";
            break;
        }
    }

    jackpot_engine_destroy(eng);
    return 0;
}
