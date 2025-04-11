#ifndef GATES_H
#define GATES_H

#include <SDL3/SDL.h>

void draw_arc(SDL_Renderer* renderer, float cx, float cy, float radius, float start_angle, float end_angle, int segments);
void draw_and_gate(SDL_Renderer* renderer, float x, float y, float w, float h);
void draw_or_gate(SDL_Renderer* renderer, float x, float y, float w, float h);
void draw_not_gate(SDL_Renderer* renderer, float x, float y, float w, float h);
void draw_xor_gate(SDL_Renderer* renderer, float x, float y, float w, float h);

#endif