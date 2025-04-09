// main.c
#include <SDL3/SDL.h>
#include "app.h"

int main(int argc, char* argv[]) {
    SDL_Init(SDL_INIT_VIDEO);
    SDL_Window* window = SDL_CreateWindow("Logic Gate Simulator", 800, 600, 0);
    SDL_Renderer* renderer = SDL_CreateRenderer(window, NULL);

    App app;
    app_init(&app, renderer); // Pass renderer to app_init

    bool running = true;
    SDL_Event event;
    while (running) {
        while (SDL_PollEvent(&event)) {
            if (event.type == SDL_EVENT_QUIT) running = false;
            app_handle_event(&app, &event);
        }
        app_update(&app);
        SDL_SetRenderDrawColor(renderer, 255, 255, 255, 255); // White background
        SDL_RenderClear(renderer);
        app_render(&app, renderer);
        SDL_RenderPresent(renderer);
    }

    app_cleanup(&app);
    SDL_DestroyRenderer(renderer);
    SDL_DestroyWindow(window);
    SDL_Quit();
    return 0;
}
