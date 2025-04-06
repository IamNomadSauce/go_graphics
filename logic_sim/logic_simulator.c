#include "logic_simulator.h"
#include <string.h>

static SDL_Renderer *g_renderer;
static GList *gates = NULL;
static GList *wires = NULL;
static Tool current_tool = TOOL_NONE;
static Gate *wire_start_gate = NULL;
static SDL_bool is_drawing_wire = SDL_FALSE;
static int mouse_x, mouse_y;

#define GATE_WIDTH 50
#define GATE_HEIGHT 30
#define CLICK_TOLERANCE 5

static Gate* find_gate_at_position(int x, int y, int *is_output, int *input_index) {
    for (GList *l = gates; l != NULL; l = l->next) {
        Gate *gate = (Gate *)l->data;
        int out_x = gate->x + gate->width + 20;
        int out_y = gate->y + 15;
        if (abs(x - out_x) < CLICK_TOLERANCE && abs(y - out_y) < CLICK_TOLERANCE) {
            *is_output = 1;
            *input_index = -1;
            return gate;
        }
        int in1_x = gate->x - 20;
        int in1_y = gate->y + 10;
        if (abs(x - in1_x) < CLICK_TOLERANCE && abs(y - in1_y) < CLICK_TOLERANCE) {
            *is_output = 0;
            *input_index = 0;
            return gate;
        }
        int in2_x = gate->x - 20;
        int in2_y = gate->y + 20;
        if (abs(x - in2_x) < CLICK_TOLERANCE && abs(y - in2_y) < CLICK_TOLERANCE) {
            *is_output = 0;
            *input_index = 1;
            return gate;
        }
    }
    return NULL;
}

static void draw_gate(Gate *gate) {
    SDL_SetRenderDrawColor(g_renderer, 0, 0, 0, 255);
    SDL_Rect rect = {gate->x, gate->y, gate->width, gate->height};
    SDL_RenderDrawRect(g_renderer, &rect);
    // Inputs
    SDL_RenderDrawLine(g_renderer, gate->x, gate->y + 10, gate->x - 20, gate->y + 10);
    SDL_RenderDrawLine(g_renderer, gate->x, gate->y + 20, gate->x - 20, gate->y + 20);
    // Output
    SDL_RenderDrawLine(g_renderer, gate->x + gate->width, gate->y + 15,
                       gate->x + gate->width + 20, gate->y + 15);
    // Label (rudimentary, SDL_ttf would be better for text)
    SDL_Point text_pos = {gate->x + 10, gate->y + 10};
    SDL_SetRenderDrawColor(g_renderer, 0, 0, 255, 255);
    for (int i = 0; gate->type[i]; i++) {
        SDL_Rect char_rect = {text_pos.x + i * 8, text_pos.y, 6, 10};
        SDL_RenderDrawRect(g_renderer, &char_rect); // Placeholder for actual text
    }
}

void logic_simulator_init(SDL_Renderer *renderer) {
    g_renderer = renderer;
    // Toolbar buttons (simulated with keyboard for now)
    SDL_Log("Use keys: 1=AND, 2=OR, 3=XOR, 4=Wire");
}

void logic_simulator_handle_event(SDL_Event *event) {
    switch (event->type) {
        case SDL_EVENT_KEY_DOWN:
            switch (event->key.keysym.sym) {
                case SDLK_1: current_tool = TOOL_AND; break;
                case SDLK_2: current_tool = TOOL_OR; break;
                case SDLK_3: current_tool = TOOL_XOR; break;
                case SDLK_4: current_tool = TOOL_WIRE; break;
            }
            break;
        case SDL_EVENT_MOUSE_BUTTON_DOWN:
            if (event->button.button == SDL_BUTTON_LEFT) {
                int x = event->button.x, y = event->button.y;
                if (current_tool == TOOL_WIRE) {
                    int is_output, input_index;
                    Gate *gate = find_gate_at_position(x, y, &is_output, &input_index);
                    if (gate) {
                        if (is_output && !is_drawing_wire) {
                            wire_start_gate = gate;
                            is_drawing_wire = SDL_TRUE;
                        } else if (!is_output && is_drawing_wire && gate != wire_start_gate) {
                            Wire *new_wire = malloc(sizeof(Wire));
                            new_wire->from = wire_start_gate;
                            new_wire->from_output = 0;
                            new_wire->to = gate;
                            new_wire->to_input = input_index;
                            wires = g_list_append(wires, new_wire);
                            is_drawing_wire = SDL_FALSE;
                        }
                    } else if (is_drawing_wire) {
                        is_drawing_wire = SDL_FALSE;
                    }
                } else if (current_tool != TOOL_NONE && current_tool != TOOL_WIRE) {
                    Gate *new_gate = malloc(sizeof(Gate));
                    new_gate->x = x;
                    new_gate->y = y;
                    new_gate->width = GATE_WIDTH;
                    new_gate->height = GATE_HEIGHT;
                    if (current_tool == TOOL_AND) new_gate->type = strdup("AND");
                    else if (current_tool == TOOL_OR) new_gate->type = strdup("OR");
                    else if (current_tool == TOOL_XOR) new_gate->type = strdup("XOR");
                    gates = g_list_append(gates, new_gate);
                    current_tool = TOOL_NONE;
                }
            }
            break;
        case SDL_EVENT_MOUSE_MOTION:
            mouse_x = event->motion.x;
            mouse_y = event->motion.y;
            break;
    }
}

void logic_simulator_draw(SDL_Renderer *renderer) {
    // Draw gates
    for (GList *l = gates; l != NULL; l = l->next) {
        draw_gate((Gate *)l->data);
    }
    // Draw wires
    SDL_SetRenderDrawColor(renderer, 0, 255, 0, 255);
    for (GList *l = wires; l != NULL; l = l->next) {
        Wire *wire = (Wire *)l->data;
        int start_x = wire->from->x + wire->from->width + 20;
        int start_y = wire->from->y + 15;
        int end_x = wire->to->x - 20;
        int end_y = wire->to->y + 10 + wire->to_input * 10;
        SDL_RenderDrawLine(renderer, start_x, start_y, end_x, end_y);
    }
    // Draw temporary wire
    if (is_drawing_wire && wire_start_gate) {
        int start_x = wire_start_gate->x + wire_start_gate->width + 20;
        int start_y = wire_start_gate->y + 15;
        SDL_RenderDrawLine(renderer, start_x, start_y, mouse_x, mouse_y);
    }
}

void logic_simulator_cleanup(void) {
    for (GList *l = gates; l != NULL; l = l->next) {
        Gate *gate = (Gate *)l->data;
        free(gate->type);
        free(gate);
    }
    g_list_free(gates);
    for (GList *l = wires; l != NULL; l = l->next) {
        free((Wire *)l->data);
    }
    g_list_free(wires);
}