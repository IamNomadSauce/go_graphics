#ifndef TYPES_H
#define TYPES_H

#include <SDL3/SDL.h>
#include <SDL3_ttf/SDL_ttf.h>
#include <glib.h>

// Forward declaration of Node
typedef struct Node Node;

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
    Node* from;  // Now valid due to forward declaration
    Node* to;
    int to_input_index;
} Wire;

typedef struct Node {
    NodeType type;
    int x, y;
    union {
        struct { bool value; } input;
        struct {
            GateType gate_type;
            GList* inputs; // List of Wire*
            bool output;
        } gate;
        struct { Wire* input; bool value; } output;
    } u;
    GList* connection_points;
} Node;

typedef struct {
    GList* nodes;
    GList* wires;
    Tool current_tool;
    bool drawing_wire;
    ConnectionPoint* wire_start_cp;
    TTF_Font* font;
    SDL_Texture* label_textures[7];
    int label_widths[7];
    int label_heights[7];
} App;

#endif