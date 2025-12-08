# jackpot_engine/Makefile
CXX      := c++

SRC      := jackpot_engine_v1.cpp
CORE     := jackpot_engine_core.hpp
HDR      := jackpot_engine_v1.h

LIBDIR   := ./
# macOS Homebrew JSON include path
JSON_INCLUDE := $(shell brew --prefix nlohmann-json)/include

CXXFLAGS := -std=c++20 -O3 -DNDEBUG -flto=auto -ffast-math \
            -Wall -Wextra -Werror -fPIC -march=native -mtune=native
NO_FAST_MATH := -fno-fast-math -fno-finite-math-only

# Detect platform
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)


# Default target
all: $(LIBDIR)/libjackpot_engine_v1_linux.so $(LIBDIR)/libjackpot_engine_v1_macos.dylib

init:
	@echo "Building for $(UNAME_S) on $(UNAME_M)..."
	brew install nlohmann-json # For macOS users; on Linux use your package manager
	# sudo apt install nlohmann-json3-dev # For Debian/Ubuntu users

# Linux shared library
$(LIBDIR)/libjackpot_engine_v1_linux.so: $(SRC) | $(LIBDIR)
	$(CXX) $(CXXFLAGS) -shared $^ -o $@

# macOS universal dylib (x86_64 + arm64)
$(LIBDIR)/libjackpot_engine_v1_macos.dylib: $(SRC) | $(LIBDIR)
	$(CXX) $(CXXFLAGS) -shared \
		-target x86_64-apple-macos11 \
		-target arm64-apple-macos11 \
		-mmacosx-version-min=11.0 \
		$^ -o $@

# Create lib dir if missing
$(LIBDIR):
	mkdir -p $(LIBDIR)

# Convenience targets
linux: $(LIBDIR)/libjackpot_engine_v1_linux.so
	@echo "Linux .so built"

macos: $(LIBDIR)/libjackpot_engine_v1_macos.dylib
	@echo "macOS universal .dylib built"

clean:
	rm -f $(LIBDIR)/*.so $(LIBDIR)/*.dylib

# Show what we built
info:
	@echo "=== Built libraries ==="
	@file $(LIBDIR)/* 2>/dev/null || echo "No libraries built yet"
	@echo
	@echo "Copy these two files into your Go module:"
	@echo "   $(LIBDIR)/libjackpot_engine_v1_linux.so"
	@echo "   $(LIBDIR)/libjackpot_engine_v1_macos.dylib"

.PHONY: all linux macos clean info

# Test simulation target
test_simulation:
	$(CXX) $(CXXFLAGS) \
		-Iv1/src \
		-I$(JSON_INCLUDE) \
		$(NO_FAST_MATH) \
		v1/tests/simulation_test_1.cpp \
		-o v1/tests/simulation_test_1

	@echo "Run simulation..."
	@cd v1/tests && ./simulation_test_1

.PHONY: test_simulation
# Test simulation target
test_jackpot_stats:
	$(CXX) $(CXXFLAGS) \
		-Iv1/src \
		-I$(JSON_INCLUDE) \
		$(NO_FAST_MATH) \
		v1/tests/jackpot_stats.cpp \
		-o v1/tests/jackpot_stats

	@echo "Run simulation..."
	@cd v1/tests && ./jackpot_stats

.PHONY: test_simulation

test_simulation_1:
	$(CXX) -std=c++20 -O3 -DNDEBUG -flto=auto \
			-Wall -Wextra -Werror -fPIC -march=native -mtune=native \
			-Iv1/src -I/opt/homebrew/include \
			v1/tests/simulation_test.cpp v1/src/jackpot_engine_v1.cpp \
			-o v1/tests/simulation_test
		@echo "=== Running 100M bet certification simulation ==="
		@v1/tests/simulation_test

.PHONY: test_simulation_1

test_drop_curve:
	$(CXX) -std=c++20 -O3 -I./ -I/opt/homebrew/include \
		tests/drop_curve_test.cpp jackpot_engine_v1.cpp \
		-o tests/drop_curve_test
	@echo "=== Запуск теста кубической кривой ==="
	@tests/drop_curve_test

.PHONY: test_drop_curve
