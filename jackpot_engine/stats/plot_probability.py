import json
import sys
import matplotlib.pyplot as plt
import numpy as np
import os


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 plot_probability.py <json_file>")
        sys.exit(1)

    json_path = sys.argv[1]

    with open(json_path, "r") as f:
        data = json.load(f)

    volatility = data["volatility"]
    step = data["step"]
    runs = data["runs"]
    intervals = data["data"]

    print(f"Loaded: volatility={volatility}, step={step}, runs={runs}")

    out_dir = "plots"
    os.makedirs(out_dir, exist_ok=True)

    for interval_section in intervals:
        iv = interval_section["interval"]
        results = interval_section["results"]

        drop_points = [r["drop_point"] for r in results]
        avg_steps = [r["avg_steps"] for r in results]

        plt.figure(figsize=(12, 6))
        plt.plot(drop_points, avg_steps, label=f"{iv['min']}..{iv['max']}", linewidth=1.5)

        plt.xlabel("DropPoint simulated")
        plt.ylabel("Probability (avg drop count per run)")
        plt.title(f"Jackpot Drop Probability — volatility={volatility} step={step}")
        plt.grid(True)
        plt.legend()

        output_path = os.path.join(out_dir, f"probability_{iv['min']}_{iv['max']}.png")
        plt.savefig(output_path, dpi=150)
        plt.close()

        print(f"Saved: {output_path}")


if __name__ == "__main__":
    main()