#include <SDL3/SDL.h>
#include <SDL3_ttf/SDL_ttf.h>
#include "app.h"
#include "render.h"

static void handle_sidebar_click(App* app, int y);
static void handle_canvas_click(App* app, int x, int y);

void app_init(App* app, SDL_Renderer* renderer) {
    
    app->nodes = NULL;
    app->wires = NULL;
    app->current_tool = TOOL_NONE;
    app->drawing_wire = false;
    app->wire_start_cp = NULL;
    app->font = NULL;

    // Initialize label_textures to NULL to avoid garbage values
    for (int i = 0; i < 7; i++) {
        app->label_textures[i] = NULL;
    }

    // Initialize SDL_ttf
    if (TTF_Init() < 0) {
        SDL_Log("TTF_Init failed: %s", SDL_GetError());
        return;
    }

    // Load the font
    app->font = TTF_OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 16);
    if (!app->font) {
        SDL_Log("Failed to load font: %s", SDL_GetError());
        return; // Exit if font loading fails
    }

    SDL_Log("Font loaded successfully");

    // Define button labels
    const char* labels[] = {"Input", "AND", "OR", "NOT", "XOR", "Wire", "Output"};
    for (int i = 0; i < 7; i++) {
        // Create surface with corrected TTF_RenderText_Blended call
        // SDL_Surface* surface = TTF_RenderText_Blended(app->font, labels[i], (SDL_Color){255, 255, 255, 255});
        SDL_Surface* surface = TTF_RenderText_Blended(
            app->font,
            labels[i],
            strlen(labels[i]),
            (SDL_Color){255, 255, 255, 255}
        );
        if (!surface) {
            SDL_Log("Failed to render text for %s: %s", labels[i], SDL_GetError());
            continue;
        }
        if (surface) {
            SDL_Log("Surface for |%s| w=%d, h=%d", labels[i], surface->w, surface->h);
        }
        SDL_Log("Surface created for '%s'", labels[i]);

        // Create texture from surface
        app->label_textures[i] = SDL_CreateTextureFromSurface(renderer, surface);
        if (!app->label_textures[i]) {
            SDL_Log("Failed to create texture for %s: %s", labels[i], SDL_GetError());
        } else {
            SDL_Log("Texture created for '%s'", labels[i]);
        }
        SDL_DestroySurface(surface); // Free surface immediately after use
    }
}

void app_handle_event(App* app, SDL_Event* event) {
    if (event->type == SDL_EVENT_MOUSE_BUTTON_DOWN) {
        int x = event->button.x;
        int y = event->button.y;
        if (x < 100) {
            handle_sidebar_click(app, y);
        } else {
            handle_canvas_click(app, x, y);
        }
    }
}

void app_update(App* app) {
    for (GList* l = app->nodes; l != NULL; l = l->next) {
        Node* node = (Node*)l->data;
        if (node->type == NODE_GATE) {
            bool input_values[2] = {false, false};
            int input_count = 0;
            for (GList* in = node->u.gate.inputs; in != NULL && input_count < 2; in = in->next) {
                Wire* wire = (Wire*)in->data;
                if (wire->from->type == NODE_INPUT) {
                    input_values[input_count++] = wire->from->u.input.value;
                } else if (wire->from->type == NODE_GATE) {
                    input_values[input_count++] = wire->from->u.gate.output;
                }
            }
            switch (node->u.gate.gate_type) {
                case GATE_AND: node->u.gate.output = input_values[0] && input_values[1]; break;
                case GATE_OR: node->u.gate.output = input_values[0] || input_values[1]; break;
                case GATE_NOT: node->u.gate.output = !input_values[0]; break;
                case GATE_XOR: node->u.gate.output = input_values[0] != input_values[1]; break;
            }
        }
    }
}

void app_render(App* app, SDL_Renderer* renderer) {
    SDL_SetRenderDrawColor(renderer, 255, 255, 255, 255);
    SDL_RenderClear(renderer);

    render_sidebar(app, renderer);
    render_nodes(app, renderer);
    render_wires(app, renderer);

    SDL_RenderPresent(renderer);
}

void app_cleanup(App* app) {
    // Free nodes and wires (existing code)
    g_list_free_full(app->nodes, free);
    g_list_free_full(app->wires, free);
    app->nodes = NULL;
    app->wires = NULL;

    // Free label textures
    for (int i = 0; i < 7; i++) {
        if (app->label_textures[i]) {
            SDL_DestroyTexture(app->label_textures[i]);  // Use SDL_DestroyTexture
            app->label_textures[i] = NULL;
        }
    }

    // Free font and TTF cleanup (if applicable)
    if (app->font) {
        TTF_CloseFont(app->font);
        app->font = NULL;
    }
    TTF_Quit();
}

static void handle_sidebar_click(App* app, int y) {
    int index = (y - 10) / 50;
    if (index >= 0 && index < 7) {
        app->current_tool = index + 2;
    }
}

static void handle_canvas_click(App* app, int x, int y) {
    if (app->current_tool == TOOL_WIRE) {
        if (!app->drawing_wire) {
            app->drawing_wire = true;
            app->wire_start_cp = malloc(sizeof(ConnectionPoint));
            app->wire_start_cp->x = x;
            app->wire_start_cp->y = y;
        } else {
            Wire* new_wire = malloc(sizeof(Wire));
            new_wire->from = NULL;
            new_wire->to = NULL;
            new_wire->to_input_index = 0;
            app->wires = g_list_append(app->wires, new_wire);
            app->drawing_wire = false;
            free(app->wire_start_cp);
            app->wire_start_cp = NULL;
        }
    } else if (app->current_tool >= TOOL_INPUT && app->current_tool <= TOOL_XOR) {
        Node* new_node = malloc(sizeof(Node));
        new_node->x = x;
        new_node->y = y;
        new_node->connection_points = NULL;
        if (app->current_tool == TOOL_INPUT) {
            new_node->type = NODE_INPUT;
            new_node->u.input.value = false; // Corrected from u.output.value
        } else if (app->current_tool == TOOL_OUTPUT) {
            new_node->type = NODE_OUTPUT;
            new_node->u.output.input = NULL;
            new_node->u.output.value = false;
        } else {
            new_node->type = NODE_GATE;
            new_node->u.gate.gate_type = app->current_tool - TOOL_AND;
            new_node->u.gate.inputs = NULL;
            new_node->u.gate.output = false;
        }
        app->nodes = g_list_append(app->nodes, new_node);
    }
}