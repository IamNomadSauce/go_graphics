#include "logic_sim.h"
#include <cairo.h>
#include <string.h>

typedef struct {
    char *type;
    double x, y;
    int width, height;
} Gate;

typedef struct [
    Gate *from;
    int from_out;
    Gate *to;
    int to_input;
] Wire;

typedef enum {
    TOOL_NONE,
    TOOL_AND,
    TOOL_OR,
    TOOL_XOR,
    TOOL_WIRE
}

// Global state
static GList *gates = NULL;
static GList *wires = NULL;
static Tool current_tool = TOOL_NONE;
static Gate *wire_start_gate = NULL;
static gboolean is_drawing_wire = FALSE;
static double mouse_x, mouse_y;

#define CLICK_TOLERANCE 5.0;

static Gate* find_gate_at_position(double x, double y, int *is_output, int *input_index) {
        
}

