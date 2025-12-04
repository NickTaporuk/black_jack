#include <iostream>
#include <fstream>
#include <vector>
#include <cstdint>
#include <stdexcept>
#include <filesystem>

#include "../src/jackpot_engine_v1.h"
#include "../src/jackpot_engine_v1.cpp"
#include "../src/jackpot_engine_core.hpp"

// Небольшая проверка ошибок
static void ensure(bool cond, const char* msg) {
    if (!cond) {
        std::cerr << msg << "\n";
        std::exit(1);
    }
}

int main() {
    std::cout << "Run jackpot stats simulation...\n";

    // 1. Инициализируем движок с каким-нибудь сидом
    uint64_t seed       = 12345678910ULL;
    uint64_t minPoint   = 1000;      // минимальный возможный dropPoint
    uint64_t maxPoint   = 2000000;   // максимальный возможный dropPoint
    int      volatility = 3;         // 1=Low, 2=Medium, 3=High (как в core)

    JackpotEngineHandle* eng = jackpot_engine_v1(seed);
    ensure(eng != nullptr, "Engine init failed");

    int id = jackpot_add(eng, minPoint, maxPoint, volatility);
    ensure(id >= 0, "jackpot_add failed");

    std::cout << "Jackpot created with id=" << id << "\n";

    // 2. Параметры симуляции
    const uint64_t TOTAL_STEPS = 5'000'000; // сколько всего шагов делаем
    const uint64_t STEP_SIZE   = 1;         // на сколько увеличиваем counter за шаг
    uint64_t lastHitStep       = 0;         // шаг, на котором был предыдущий джекпот
    uint64_t hitCount          = 0;         // количество срабатываний


    namespace fs = std::filesystem;
    fs::path base = fs::current_path();
    std::cout << base <<std::endl;
    fs::path out  = base.parent_path() / "stats" / "jackpot_stats.csv";

    // Создаём директорию stats, если её ещё нет
    fs::create_directories(out.parent_path());

    std::cout << "CSV will be saved to " << out << "\n";

    std::ofstream csv(out);
    if (!csv.is_open()) {
        std::cerr << "Failed to open " << out << " for writing\n";
        jackpot_engine_destroy(eng);
        return 1;
    }

    if (!csv.is_open()) {
        std::cerr << "Failed to open jackpot_stats.csv for writing\n";
        jackpot_engine_destroy(eng);
        return 1;
    }

    // Заголовки
    csv << "hit_index,hit_distance,drop_point\n";

    // 4. Основной цикл
    for (uint64_t step = 1; step <= TOTAL_STEPS; ++step) {
        // dropPoint до шага — это тот, который сейчас активен
        uint64_t dp_before = jackpot_get_drop_point(eng, id);

        int hit = jackpot_process(eng, id, STEP_SIZE);

        if (hit == 1) {
            ++hitCount;
            uint64_t distance = step - lastHitStep;
            lastHitStep = step;

            // Логируем в CSV
            csv << hitCount << "," << distance << "," << dp_before << "\n";

            if (hitCount % 10 == 0) {
                std::cout << "Hit #" << hitCount
                          << " distance=" << distance
                          << " dropPoint(before)=" << dp_before
                          << " currentCounter=" << jackpot_get_counter(eng, id)
                          << " newDropPoint=" << jackpot_get_drop_point(eng, id)
                          << "\n";
            }
        }
    }

    std::cout << "Simulation finished. Total hits = " << hitCount << "\n";
    std::cout << "CSV saved to v1/tests/jackpot_stats.csv\n";

    jackpot_engine_destroy(eng);
    csv.close();
    return 0;
}