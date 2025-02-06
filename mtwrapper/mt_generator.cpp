#include <iostream>
#include <random>
#include <ctime>
#include <cstdlib>  // For malloc and free

extern "C" {
    // Properly seeded 64-bit Mersenne Twister
    std::mt19937_64 mt_engine(static_cast<unsigned long>(std::time(nullptr)));

    // Generate random tokens
    uint64_t* generate_random_tokens(uint64_t rangeGen) {
        uint64_t* tokens = (uint64_t*)malloc(rangeGen * sizeof(uint64_t)); // Allocate memory

        std::uniform_int_distribution<uint64_t> dist(1, UINT64_MAX); // Avoid zero

        for (uint64_t i = 0; i < rangeGen; i++) {
            tokens[i] = dist(mt_engine);
        }

        return tokens;
    }

    // Get a random index from the tokens
    int get_random_index(uint64_t* tokens, uint64_t rangeGen) {
        std::uniform_int_distribution<int> dist(0, rangeGen - 1);
        return dist(mt_engine);
    }

    // Get a random number from the tokens
    uint64_t get_random_number(uint64_t* tokens, uint64_t rangeGen) {
        int randomIndex = get_random_index(tokens, rangeGen);
        return tokens[randomIndex];
    }

    // Drop bonus function
    bool drop_bonus(uint64_t rangeGen) {
        uint64_t* tokens = generate_random_tokens(rangeGen);
        uint64_t randomNumber1 = get_random_number(tokens, rangeGen);
        uint64_t randomNumber2 = get_random_number(tokens, rangeGen);

        bool result = (randomNumber1 == randomNumber2);
        free(tokens);  // Free allocated memory

        return result;
    }
}