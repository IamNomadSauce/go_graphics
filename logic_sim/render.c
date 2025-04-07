#include "render.h"

void render_sidebar(App* app, SDL_Renderer renderer) {
    SDL_SetRenderDrawColor(renderer, 200, 200, 200, 255);
    SDL_Rect sidebar = {0,0,100, 600};
    SDL_RenderFillRect(renderer, &sidebar);

    const char* labels[] = {"Input", "AND", "OR", "NOT", "XOR", "Wire", "Output"};
    for (int i = 0; i < 7; i++) {
        SDL_Rect button = {10, 10 + 1 * 50, 80, 40};
        if (app->current_tool == i + 1) {
            SDL_SetRenderDrawColor(renderer, 150, 150, 255, 255);
        } else {
            SDL_SetRenderDrawColor(render, 100, 100, 100, 255);
        }
        SDL_RenderFillRect(renderer, &button);
    }
}

void render_nodes(App* app, SDL_Renderer* renderer) {

}

void render_wires(App* app, SDL_Rednerer* renderer) {

}