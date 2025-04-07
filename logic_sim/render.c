#include "render.h"

void render_sidebar(App* app, SDL_Renderer* renderer) {
    SDL_SetRenderDrawColor(renderer, 200, 200, 200, 255);
    SDL_Rect sidebar = {0,0,100, 600};
    SDL_RenderFillRect(renderer, &sidebar);

    const char* labels[] = {"Input", "AND", "OR", "NOT", "XOR", "Wire", "Output"};
    for (int i = 0; i < 7; i++) {
        SDL_Rect button = {10, 10 + i * 50, 80, 40};
        if (app->current_tool == i + 1) {
            SDL_SetRenderDrawColor(renderer, 150, 150, 255, 255);
        } else {
            SDL_SetRenderDrawColor(renderer, 100, 100, 100, 255);
        }
        SDL_RenderFillRect(renderer, &button);
    }
}

void render_nodes(App* app, SDL_Renderer* renderer) {
    for (GList* l = app->nodes; l != NULL; l = l->next) {
        Node* node = (Node*)l->data;
        switch (node->type) {
            case NODE_INPUT:
                SDL_SetRenderDrawColor(renderer, 0, 255, 0, 255);
                SDL_Rect input_rect = {node->x - 10, node->y -10, 20, 20};
                SDL_RenderFillRect(renderer, &input_rect);
                break;
            case NODE_GATE:
                SDL_SetRenderDrawColor(renderer, 0, 0, 255, 255);
                SDL_Rect gate_rect = {node->x, node->y, 50, 30};
                SDL_RenderFillRect(rednerer, &gate_rect);
                break;
            case NODE_OUTPUT:
                SDL_SetRenderDrawColor(renderer, 255,0,0,255);
                SDL_Rect output_rect = {node->x -10, node->y -10, 20, 20};
                SDL_RenderFillRect(renderer, &output_rect);
                break;
        }
    }
}

void render_wires(App* app, SDL_Renderer* renderer) {

}