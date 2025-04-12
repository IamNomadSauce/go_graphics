#include <SDL3/SDL.h>
#include "gates.h"
#include <math.h>

void draw_arc(SDL_Renderer* renderer, float cx, float cy, float radius, float start_angle, float end_angle, int segments) {
    float angle_step = (end_angle - start_angle) / segments;
    for (int i = 0; i < segments; i++) {
        float theta1 = start_angle + i * angle_step;
        float theta2 = start_angle + (i + 1) * angle_step;
        float x1 = cx + radius * cos(theta1);
        float y1 = cy + radius * sin(theta1);
        float x2 = cx + radius * cos(theta2);
        float y2 = cy + radius * sin(theta2);
        SDL_RenderLine(renderer, x1, y1, x2, y2);
    }
}

void draw_and_gate(SDL_Renderer* renderer, float x, float y, float w, float h) {
    // SDL_Log("AND Gate");
    const int arc_segments = 10;
    const float r = h / 2.0f;
    const float cx = x + w - r;
    const float cy = y + h / 2.0f;

    SDL_RenderLine(renderer, x, y, x, y + h);

    SDL_RenderLine(renderer, x, y, x + w - r, y);

    draw_arc(renderer, cx, cy, r, -M_PI / 2, M_PI / 2, arc_segments);

    SDL_RenderLine(renderer, x + w - r, y + h, x, y + h);
}
    
void draw_or_gate(SDL_Renderer* renderer, float x, float y, float w, float h) {
    // SDL_Log("OR Gate");
    const int arc_segments = 10;
    const float input_radius = w / 2.0f;
    const float input_cx = x - w / 4.0f;
    const float cy = y + h / 2.0f;


    draw_arc(renderer, input_cx, cy, input_radius, M_PI / 3, 2 * M_PI / 3, arc_segments);

    SDL_RenderLine(renderer, x + w / 4, y, x + w, y + h /2);

    SDL_RenderLine(renderer, x + w / 4, y + h, x + w, y + h / 2);
}

void draw_not_gate(SDL_Renderer* renderer, float x, float y, float w, float h) {
    // SDL_Log("NOT Gate");
    const int circle_segments = 10;
    const float circle_radius = 5.0f;
    const float circle_cx = x + w - circle_radius;
    const float circle_cy = y + h / 2.0f;

    SDL_RenderLine(renderer, x, y + h / 2, x + w - 2 * circle_radius, y);
    SDL_RenderLine(renderer, x + w - 2 * circle_radius, y, x + w - 2 * circle_radius, y + h);
    SDL_RenderLine(renderer, x + w - 2 * circle_radius, y + h, x, y + h / 2);

    draw_arc(renderer, circle_cx, circle_cy, circle_radius, 0, 2 * M_PI, circle_segments);
}

void draw_xor_gate(SDL_Renderer* renderer, float x, float y, float w, float h) {
    // SDL_Log("XOR Gate");
    const int arc_segments = 10;
    const float input_radius = w / 2.0f;
    const float input_cx = x - w / 4.0f;
    const float cy = y + h / 2.0f;

    draw_arc(renderer, input_cx, cy, input_radius, M_PI / 3, 2 * M_PI / 3, arc_segments);
    SDL_RenderLine(renderer, x + w / 4, y, x + w, y + h / 2);
    SDL_RenderLine(renderer, x + w / 4, y, x + w, y + h /2);
    SDL_RenderLine(renderer, x + w / 4, y + h, x + w,  y + h / 2);

    const float xor_arc_radius = input_radius * 0.8f;
    const float xor_arc_cx = input_cx - w/ 8.0f;
    draw_arc(renderer, xor_arc_cx, cy, xor_arc_radius, M_PI / 3, 2 * M_PI / 3, arc_segments);
}