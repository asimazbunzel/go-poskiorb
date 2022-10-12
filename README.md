# Binary orbits after asymmetric momentum kick

Compute binary orbits after an asymmetry in a core collapse of a massive star using linear
momentum conservation as presented in Kalogera et al. 1996 and Kalogera 2000

## Installation

The code is written purely in Go so its installation requires `go` to be installed in your system.

First, clone this code into your computer

```
git clone git@github.com:asimazbunzel/go-orbits.git
```

(in case ssh is not available, clone it with the HTTPS URL).

Then, the module must be initialized with

```
go mod init github.com/asimazbunzel/go-orbits
```

Once this is done, the following command will get you all the libraries needed for the code to
work

```
go mod tidy
```

Finally, to compile the code run

```
make build
```

(which needs the command `make` available, always present in a GNU/Linux distribution).

This last command will generate a binary called `orbits` located in  the `bin` folder. This is
the program that will be called to compute the distribution of kicks.

## Usage

All the necessary options are controlled via a configuration file in YAML format. There is an
example in this repository called `config.yaml`. In it, there is an explanation for the utility
of each of the options.

In order to run the code, change directory into the `bin` folder and run

```
./orbits -C path/to/config/file
```

For example, in the example given with the `config.yaml` file, the command would look like this

```
./orbits -C ../config.yaml
```

### Some comments about the options

* `m1`, `m2`, `separation` and `period` are the conditions of the binary just before the core
collapse.

* `compact_object_mass` is the mass of the newly born compact object.

* `kick_distribution` and `kick_direction` are the distributions for the asymmetric kick. The
`kick_distribution` options are either `Maxwell` or `Uniform` for a Maxwellian or a Uniform
distribution of the strength of the kick. While `kick_direction` has not other available option
as of now (maybe in the future could be changed).

* `reduce_by_fallback` applies a function that reduces the strength of the kick by `(1 - f)` with
`f` being the fraction of mass that falls back to the compact object at core collapse.

* `kick_sigma` is used when `kick_distribution` is `Maxwell`. While `min_kick_value` and
`max_kick_value` are used when `kick_distribution` is `Uniform`.

* `min_phi` and `max_phi` are options for the angle in the orbital plane of the binary. Do not
change it, for now.

* `seed` is the number used by the random number generator method.

* `number_of_cases` represents the number of draws for the different kicks.

* `log_level` option for the amount of terminal output. Options are `debug` or `info`.

* `save_kicks` and `kicks_filename` are self explanatory.

* `save_bounded_orbits` and `bounded_orbits_filename` are used to store the binaries that
survive the kick.

* `save_grid_of_orbits` and `grid_of_orbits_filename` are used to store a grid of binaries
with a probability above a threshold of `minimum_probability_for_grid`.

* `period_quantile_min`, `period_quantile_max`, `eccentricity_quantile_min`,
`eccentricity_quantile_max`, `number_of_periods` and `number_of_eccentricities` are controls
for the creation of the grid. Do not change them, for now.

## Output

The code will create 3 different files (according to some controls shown above). One of the
files will contain info on the strength and direction of the kick (`kicks_filename`), another
will have info on the binaries that survive the kick (`bounded_orbits_filename`). The last
file will create a grid of orbital parameters assuming that the 2D plane of
(period, eccentricity) can be divided into a rectangular grid in which, each of the rectangles
will have associated a probability according to how many binaries are within its boundaries
(`grid_of_orbits_filename`).
