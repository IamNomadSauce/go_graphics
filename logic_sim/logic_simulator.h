#ifndef LOGIC_SIMULATOR_H
#define LOGIC_SIMULATOR_H

#include <SDL3/SDL.h>

typedef enum {
    TOOL_NONE,
    TOOL_AND,
    TOOL_OR,
    TOOL_XOR,
    TOOL_WIRE,
} Tool;


typedef struct {
    char *type;
    int x, y;
    int width, height;
} Gate;

typedef struct {
    Gate *from;
    int from_output;
    Gate *to;
    int to_input;
} Wire;

void logic_simulator_draw(SDL_Renderer *renderer);

void logic_simulator_init(SDL_Renderer *renderer);
void logic_simulator_handle_event(SDL_Event *event);
void logc_simulator_draw(SDL_Renderer *renderer);
void logic_simulator_cleanup(void);

#endif