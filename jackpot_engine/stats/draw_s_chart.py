import numpy as np
import matplotlib.pyplot as plt

x = np.linspace(0, 1, 1000)
p = 10*x**5 - 15*x**4 + 6*x**3

plt.figure(figsize=(10, 6), dpi=300)
plt.plot(x, p, label='p(x) = 10x⁵ − 15x⁴ + 6x³', linewidth=4, color='#d62728')
plt.plot(x, x, '--', label='Линейная (Low vol)', alpha=0.6)
plt.plot(x, x**0.5, '--', label='√x (Medium vol)', alpha=0.6)
plt.plot(x, x**3, '--', label='x³ (High vol)', alpha=0.6)

plt.xlabel('Прогресс (counter / dropPoint)')
plt.ylabel('Мгновенная вероятность выпадения ×10,000')
plt.title('Твоя квинтическая кривая vs стандартные волатильности')
plt.legend()
plt.grid(True, alpha=0.3)
plt.tight_layout()
plt.savefig('quintic_volatility_curve.pdf')
plt.savefig('quintic_volatility_curve.png', dpi=300)
plt.show()