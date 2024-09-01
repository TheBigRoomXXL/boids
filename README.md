# Boids 

This program simulates "boids," based on the influential paper [Flocks, Herds, and Schools:
A Distributed Behavioral Model](http://www.red3d.com/cwr/papers/1987/boids.html) . The implementation follows the pseudo code provided by [V. Hunter Adams](https://vanhunteradams.com/Pico/Animal_Movement/Boids-algorithm.html).

The program use the awesome [raylite](https://www.raylib.com/) library for rendering.

![demo](demo.gif)

## Thing I might (or might not) do in the future

- Implement a more efficient algorithm to simulate a larger number of boids.
- Add support for user inputs via the CLI and directly within the rendering interface (e.g., sliders for adjusting factors).
- Target HTML to enable the simulation to run directly in a web browser.
- Introduce more complexe behaviors such as predator-prey dynamics or multiple boid groups.
