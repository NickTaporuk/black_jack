package v2

/*
#cgo  darwin LDFLAGS: -Wl,-rpath,"${SRCDIR}"
#cgo  linux LDFLAGS: -L${SRCDIR}"

#include <stdint.h>
#include "jackpot_c_api.h"
*/
import "C"
import "errors"

func CheckUniform(
	seed, min, max, current uint64,
) (dropped bool, dropPoint uint64, err error) {

	var dp C.uint64_t

	ok := C.jackpot_check_uniform(
		C.uint64_t(seed),
		C.uint64_t(min),
		C.uint64_t(max),
		C.uint64_t(current),
		&dp,
	)

	if min >= max {
		return false, 0, errors.New("invalid interval")
	}

	return ok == 1, uint64(dp), nil
}
