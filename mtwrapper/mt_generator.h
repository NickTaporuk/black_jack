#ifndef MT_GENERATOR_H
#define MT_GENERATOR_H

#include <stdint.h>   // uint64_t
#include <stdbool.h>  // bool
#include <stdlib.h>   // free()

#ifdef __cplusplus
extern "C" {
#endif

/**
 * Generates an array of random uint64_t values.
 * IMPORTANT: Memory is allocated using malloc() — the caller MUST free()
 * the returned pointer when it is no longer needed.
 *
 * @param rangeGen  the number of tokens to generate
 * @return a pointer to an array of uint64_t values (or NULL on failure)
 */
uint64_t* generate_random_tokens(uint64_t rangeGen);

/**
 * Returns a random index in the range [0, rangeGen - 1].
 *
 * @param tokens     pointer to an array of uint64_t values
 * @param rangeGen   number of elements in the tokens array
 * @return a random index
 */
int get_random_index(uint64_t* tokens, uint64_t rangeGen);

/**
 * Returns a random uint64_t value from the provided array.
 *
 * @param tokens     pointer to an array of uint64_t values
 * @param rangeGen   number of elements in the tokens array
 * @return a randomly selected value from the array
 */
uint64_t get_random_number(uint64_t* tokens, uint64_t rangeGen);

/**
 * Generates two random numbers and checks whether they are equal.
 * Effectively models a "bonus drop" event where a match results in success.
 *
 * @param rangeGen   number of tokens used in the internal generation
 * @return true if the two generated values match, false otherwise
 */
bool drop_bonus(uint64_t rangeGen);

#ifdef __cplusplus
}
#endif

#endif // MT_GENERATOR_H
