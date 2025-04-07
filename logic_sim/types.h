#ifndef TYPES_H
#define TYPES_H

#include <SDL3/SDL.h>
#include <glib.h>

typedef enum {
    NODE_INPUT,
    NODE_GATE,
    NODE_OUTPUT
} NodeType;

typedef enum {
    GATE_AND,
    GATE_OR,
    GATE_NOT,
    GATE_XOR
} GateType;

typedef enum {
    TOOL_NONE,
    TOOL_WIRE,
    TOOL_INPUT,
    TOOL_OUTPUT,
    TOOL_AND,
    TOOL_OR,
    TOOL_NOT,
    TOOL_XOR
} Tool;

typedef struct {
    int x, y;
    bool is_input;
    int index;
} ConnectionPoint;

typedef struct {
    NodeType type;
    int x, y;
    union {
        struct { bool value; } input;
        struct {
            GateType gate_type;
            GList* inputs;
            bool output;
        } gate;
        struct { wire* input; bool value; } output;
    } u;
    GList* connection_points;
} Node;

typedef struct {
    Node* from;
    Node* to;
    int to_input_index;
} Wire;

typedef struct {
    GList* nodes;
    GList* wires;
    Tool current_tool;
    bool drawing_wire;
    ConnectionPoint* wire_start_cp;
} App;

#endif