
import matplotlib.pyplot as plt
import numpy as np


# first, plot kick distribution between this module and a python one
index_g, w_g, theta_g, phi_g = np.loadtxt("kicks.data", skiprows=1, unpack=True)

# kick strength
fig, ax = plt.subplots()
ax.hist(w_g, bins=50, histtype="step", color="C0")
plt.show()

# azimuthal angle
fig, ax = plt.subplots()
ax.hist(theta_g, bins=50, histtype="step", color="C0")
plt.show()

# polar angle
fig, ax = plt.subplots()
ax.hist(phi_g, bins=50, histtype="step", color="C0")
plt.show()


# load and compare orbit distributions
index_g, _, _, _, p_g, a_g, e_g = np.loadtxt("orbits.data", skiprows=1, unpack=True)

fig, ax = plt.subplots()
ax.set_xscale("log")
ax.set_xlim(3, 400)
ax.set_ylim(0, 1)
ax.scatter(p_g, e_g, s=1, c="C0")
plt.show()


# load and compare grid of orbits
index_g, p_g, a_g, e_g, prob_g = np.loadtxt("grid.data", skiprows=1, unpack=True)

fig, ax = plt.subplots()
ax.set_xscale("log")
ax.set_xlim(3, 400)
ax.set_ylim(0, 1)
ax.scatter(p_g, e_g, s=1, c="C0")
plt.show()
