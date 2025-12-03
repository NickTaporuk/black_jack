import csv
import sys
import matplotlib.pyplot as plt

def read_csv(path):
    hit_idx = []
    hit_distance = []
    drop_points = []

    with open(path, "r") as f:
        reader = csv.DictReader(f)
        for row in reader:
            hit_idx.append(int(row["hit_index"]))
            hit_distance.append(int(row["hit_distance"]))
            drop_points.append(int(row["drop_point"]))

    return hit_idx, hit_distance, drop_points


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 plot_report.py jackpot_stats.csv")
        return

    csv_path = sys.argv[1]
    print(f"Reading CSV: {csv_path}")

    hit_idx, distances, drop_points = read_csv(csv_path)

    # ---------------------------------------------
    # 1) Scatter plot: hit_distance
    # ---------------------------------------------
    plt.figure(figsize=(12, 6))
    plt.scatter(hit_idx, distances, s=10, alpha=0.6, label="Hit Distance")
    plt.title("Jackpot Hit Distance Over Time")
    plt.xlabel("Hit index")
    plt.ylabel("Hit distance")
    plt.grid(True)
    plt.legend()
    plt.tight_layout()
    plt.savefig("jackpot_hit_distance.png", dpi=200)
    print("Saved jackpot_hit_distance.png")

    # ---------------------------------------------
    # 2) Line plot: drop_point evolution
    # ---------------------------------------------
    plt.figure(figsize=(12, 6))
    plt.plot(hit_idx, drop_points, linewidth=1.3, label="DropPoint")
    plt.title("Drop Point evolution")
    plt.xlabel("Hit index")
    plt.ylabel("Drop Point value")
    plt.grid(True)
    plt.legend()
    plt.tight_layout()
    plt.savefig("jackpot_drop_point.png", dpi=200)
    print("Saved jackpot_drop_point.png")

    # ---------------------------------------------
    # 3) Scatter vs Expected line
    # ---------------------------------------------
    expected = [dp * 0.5 for dp in drop_points]  # пример модели ожидания

    plt.figure(figsize=(12, 6))
    plt.scatter(hit_idx, distances, s=10, alpha=0.4, label="Actual Hit Distance")
    plt.plot(hit_idx, expected, "r--", label="Expected trend (demo)")
    plt.title("Hit Distance vs Trend")
    plt.xlabel("Hit index")
    plt.ylabel("Value")
    plt.grid(True)
    plt.legend()
    plt.tight_layout()
    plt.savefig("jackpot_distance_vs_expected.png", dpi=200)
    print("Saved jackpot_distance_vs_expected.png")

    print("\n✔️ Graphs generated successfully!")


if __name__ == "__main__":
    main()