#include <SDL3/SDL.h>
#include <SDL3_ttf/SDL_ttf.h>
#include "app.h"
#include "render.h"

const Tool sidebar_tools[7] = {
    TOOL_INPUT,
    TOOL_AND,
    TOOL_OR,
    TOOL_NOT,
    TOOL_XOR,
    TOOL_WIRE,
    TOOL_OUTPUT
};
static void handle_sidebar_click(App* app, int y);
static void handle_canvas_click(App* app, int x, int y);


void app_init(App* app, SDL_Renderer* renderer) {

    SDL_Log("Renderer in app_init: %p", renderer);
    
    app->nodes = NULL;
    app->wires = NULL;
    app->current_tool = TOOL_NONE;
    app->drawing_wire = false;
    app->wire_start_cp = NULL;
    app->font = NULL;
    memset(app->label_widths, 0, sizeof(app->label_widths));
    memset(app->label_heights, 0, sizeof(app->label_heights));

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
        app->label_textures[i] = SDL_CreateTextureFromSurface(renderer, surface);
        if (!app->label_textures[i]) {
            SDL_Log("Failed to create texture for %s: %s", labels[i], SDL_GetError());
        } else {
            SDL_Log("Texture created for '%s' with w=%d, h=%d", labels[i], surface->w, surface->h);
            app->label_widths[i] = surface->w;
            app->label_heights[i] = surface->h;
        }
        SDL_DestroySurface(surface);
    }
}

void app_handle_event(App* app, SDL_Event* event) {
    // SDL_Log("app_handle_event %c", event->type);
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
    SDL_Log("handle_sidebar_click: |%d|", y);
    SDL_Log(":current_tool |%d|", app->current_tool);
    int index = (y - 10) / 50;
    if (index >= 0 && index < 7) {
        app->current_tool = sidebar_tools[index];
        SDL_Log("Index |%d| current_tool |%d|", index, app->current_tool);
    }
}

static void handle_canvas_click(App* app, int x, int y) {
    SDL_Log("handle_canvas_click: |%d| |%d|", x, y);
    SDL_Log("%d", app->current_tool);
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
        SDL_Log("Added node at (%d, %d) of type %d", x, y, new_node->type);
    }
}