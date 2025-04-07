#ifndef RENDER_H
#define RENDER_H

#include "types.h"

void render_sidebar(App* app, SDL_Renderer* renderer);
void render_nodes(App* app, SDL_Renderer* renderer);
void render_wires(App* app, SDL_Renderer* renderer);

#endif