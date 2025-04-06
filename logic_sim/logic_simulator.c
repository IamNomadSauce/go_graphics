#include "logic_simulator.h"
#include <stdlib.h>
#include <glib.h>
#include <string.h>

static SDL_Renderer *g_renderer;
static GList *gates = NULL;
static GList *wires = NULL;
static Tool current_tool = TOOL_NONE;
static Gate *wire_start_gate = NULL;
static int is_drawing_wire = 0; // Using int instead of SDL_bool
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
    SDL_SetRenderDrawColor(g_renderer, 0, 0, 0, 255); // Set color (unchanged)
    SDL_FRect frect = {(float)gate->x, (float)gate->y, (float)gate->width, (float)gate->height};
    SDL_RenderRect(g_renderer, &frect); // Updated from SDL_RenderDrawRect

    // Draw input lines
    SDL_RenderLine(g_renderer, (float)gate->x, (float)(gate->y + 10), 
                   (float)(gate->x - 20), (float)(gate->y + 10));
    SDL_RenderLine(g_renderer, (float)gate->x, (float)(gate->y + 20), 
                   (float)(gate->x - 20), (float)(gate->y + 20));

    // Draw output line
    SDL_RenderLine(g_renderer, (float)(gate->x + gate->width), (float)(gate->y + 15),
                   (float)(gate->x + gate->width + 20), (float)(gate->y + 15));

    // Example label drawing (adjust as needed)
    SDL_Point text_pos = {gate->x + 10, gate->y + 10};
    SDL_SetRenderDrawColor(g_renderer, 0, 0, 255, 255);
    for (int i = 0; gate->type[i]; i++) {
        SDL_FRect char_rect = {(float)(text_pos.x + i * 8), (float)text_pos.y, 6.0f, 10.0f};
        SDL_RenderRect(g_renderer, &char_rect); // Updated from SDL_RenderDrawRect
    }
}

void logic_simulator_init(SDL_Renderer *renderer) {
    g_renderer = renderer;
    SDL_Log("Use keys: 1=AND, 2=OR, 3=XOR, 4=Wire");
}

void logic_simulator_handle_event(SDL_Event *event) {
    switch (event->type) {
        case SDL_EVENT_KEY_DOWN:
            switch (event->key.scancode) {
                case SDL_SCANCODE_1: current_tool = TOOL_AND; break;
                case SDL_SCANCODE_2: current_tool = TOOL_OR; break;
                case SDL_SCANCODE_3: current_tool = TOOL_XOR; break;
                case SDL_SCANCODE_4: current_tool = TOOL_WIRE; break;
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
                            is_drawing_wire = 1;
                        } else if (!is_output && is_drawing_wire && gate != wire_start_gate) {
                            Wire *new_wire = malloc(sizeof(Wire));
                            new_wire->from = wire_start_gate;
                            new_wire->from_output = 0;
                            new_wire->to = gate;
                            new_wire->to_input = input_index;
                            wires = g_list_append(wires, new_wire);
                            is_drawing_wire = 0;
                        }
                    } else if (is_drawing_wire) {
                        is_drawing_wire = 0;
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
    for (GList *l = gates; l != NULL; l = l->next) {
        draw_gate((Gate *)l->data);
    }
    SDL_SetRenderDrawColor(renderer, 0, 255, 0, 255); // Green wires
    for (GList *l = wires; l != NULL; l = l->next) {
        Wire *wire = (Wire *)l->data;
        float start_x = (float)(wire->from->x + wire->from->width + 20);
        float start_y = (float)(wire->from->y + 15);
        float end_x = (float)(wire->to->x - 20);
        float end_y = (float)(wire->to->y + 10 + wire->to_input * 10);
        SDL_RenderLine(renderer, start_x, start_y, end_x, end_y); // Updated
    }
    if (is_drawing_wire && wire_start_gate) {
        float start_x = (float)(wire_start_gate->x + wire_start_gate->width + 20);
        float start_y = (float)(wire_start_gate->y + 15);
        SDL_RenderLine(renderer, start_x, start_y, (float)mouse_x, (float)mouse_y); // Updated
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