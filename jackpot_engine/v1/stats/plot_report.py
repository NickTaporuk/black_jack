import json
import matplotlib.pyplot as plt

simulation_report_path = "./../tests/simulation_report.json"
with open("simulation_report.json") as f:
    data = json.load(f)

hits = data["hits"]

progress = [h["progress"] for h in hits]

plt.hist(progress, bins=50)
plt.title("Jackpot Drop Progress Distribution")
plt.xlabel("progress = counter / dropPoint")
plt.ylabel("hits")
plt.savefig("progress_hist.png", dpi=160)
print("Saved progress_hist.png")
