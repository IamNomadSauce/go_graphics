#include <SDL3/SDL.h>
#include "render.h"
#include "gates.h"
#include <math.h>



void render_sidebar(App* app, SDL_Renderer* renderer) {
    // Draw the sidebar background
    SDL_Log("Renderer in render_sidebar: %p", renderer);
    SDL_SetRenderDrawColor(renderer, 200, 200, 200, 255);
    SDL_FRect sidebar = {0.0f, 0.0f, 100.0f, 600.0f};
    SDL_RenderFillRect(renderer, &sidebar);

    // Render each button in the sidebar
    for (int i = 0; i < 7; i++) {
        SDL_FRect button = {10.0f, 10.0f + i * 50.0f, 80.0f, 40.0f};
        if (app->current_tool == i + 2) {
            SDL_SetRenderDrawColor(renderer, 150, 150, 255, 255); // Highlight selected tool
        } else {
            SDL_SetRenderDrawColor(renderer, 100, 100, 100, 255); // Default button color
        }
        SDL_RenderFillRect(renderer, &button);

        // Render the label texture if it exists and has valid size
        if (app->label_textures[i] && app->label_widths[i] > 0 && app->label_heights[i] > 0) {
            float tex_w = (float)app->label_widths[i];
            float tex_h = (float)app->label_heights[i];
            SDL_Log("Button %d texture: w=%f, h=%f", i, tex_w, tex_h);

            SDL_FRect dest_rect = {
                .x = button.x + (button.w - tex_w) / 2.0f, // Center horizontally
                .y = button.y + (button.h - tex_h) / 2.0f, // Center vertically
                .w = tex_w,
                .h = tex_h
            };
            SDL_RenderTexture(renderer, app->label_textures[i], NULL, &dest_rect);
        } else {
            SDL_Log("No valid texture or size for button %d", i);
        }
    }
}

void render_nodes(App* app, SDL_Renderer* renderer) {
    for (GList* l = app->nodes; l != NULL; l = l->next) {
        Node* node = (Node*)l->data;
        SDL_SetRenderDrawColor(renderer, 255, 255, 255, 255);
        if (node->type == NODE_GATE) {
            switch (node->u.gate.gate_type) {
                case GATE_AND:
                    draw_and_gate(renderer, node->x, node->y, 50, 30);
                    break;
                case GATE_OR:
                    draw_or_gate(renderer, node->x, node->y, 50, 30);
                    break;
                case GATE_NOT:
                    draw_not_gate(renderer, node->x, node->y, 50, 30);
                    break;
                case GATE_XOR:
                    draw_xor_gate(renderer, node->x, node->y, 50, 30);
                    break;
            }
        } else if (node->type == NODE_INPUT) {
            draw_arc(renderer, (float)(node->x + 10), (float)(node->y + 10), 10.0f, 0.0f, 2.0f * (float)M_PI, 20);
        } else if (node->type == NODE_OUTPUT) {
            SDL_FRect rect = {node->x, node->y, 20, 20};
            SDL_RenderRect(renderer, &rect);  // Updated to SDL_RenderRect
        }
    }
}

void render_wires(App* app, SDL_Renderer* renderer) {
    for (GList* l = app->wires; l != NULL; l = l->next) {
        Wire* wire = (Wire*)l->data;
        float from_x = (float)(wire->from->x + (wire->from->type == NODE_GATE ? 50 : 10));
        float from_y = (float)(wire->from->y + 15);
        float to_x = (float)wire->to->x;
        float to_y = (float)(wire->to->y + 15 + wire->to_input_index * 10);
        SDL_SetRenderDrawColor(renderer, 0, 0, 0, 255); // Black for wires
        SDL_RenderLine(renderer, from_x, from_y, to_x, to_y);
    }
}