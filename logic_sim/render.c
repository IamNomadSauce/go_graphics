#include "render.h"

void render_sidebar(App* app, SDL_Renderer* renderer) {
    SDL_SetRenderDrawColor(renderer, 200, 200, 200, 255);
    SDL_FRect sidebar = {0.0f, 0.0f, 100.0f, 600.0f};
    SDL_RenderFillRect(renderer, &sidebar);

    const char* labels[] = {"Input", "AND", "OR", "NOT", "XOR", "Wire", "Output"};
    for (int i = 0; i < 7; i++) {
        SDL_FRect button = {10.0f, 10.0f + i * 50.0f, 80.0f, 40.0f};
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
                SDL_FRect input_rect = {(float)node->x - 10.0f, (float)node->y - 10.0f, 20.0f, 20.0f};
                SDL_RenderFillRect(renderer, &input_rect);
                break;
            case NODE_GATE:
                SDL_SetRenderDrawColor(renderer, 0, 0, 255, 255);
                SDL_FRect gate_rect = {(float)node->x, (float)node->y, 50.0f, 30.0f};
                SDL_RenderFillRect(renderer, &gate_rect);  // Fixed typo 'rednerer'
                break;
            case NODE_OUTPUT:
                SDL_SetRenderDrawColor(renderer, 255, 0, 0, 255);
                SDL_FRect output_rect = {(float)node->x - 10.0f, (float)node->y - 10.0f, 20.0f, 20.0f};
                SDL_RenderFillRect(renderer, &output_rect);
                break;
        }
    }
}

void render_wires(App* app, SDL_Renderer* renderer) {
    for (GList* l = app->wires; l != NULL; l = l->next) {
        Wire* wire = (Wire*)l->data;
        float from_x = (float)(wire->from->x + (wire->from->type == NODE_GATE ? 50 : 10));
        float from_y = (float)(wire->from->y + 15);
        float to_x = (float)wire->to->x;
        float to_y = (float)(wire->to->y + 15 + wire->to_input_index * 10);  // Fixed syntax
        SDL_SetRenderDrawColor(renderer, 0, 0, 0, 255);  // Fixed function name
        SDL_RenderLine(renderer, from_x, from_y, to_x, to_y);  // SDL3 uses SDL_RenderLine
    }
}