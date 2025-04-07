#include "app.h"
#include "render.h"

void app_init(App* app) {
    app->nodes = NULL;
    app->wires = NULL;
    app->current_tool = TOOL_NONE;
    app->drawing_wire = false;
    app->wire_start_cp = NULL;
}

void app_handle_event(App* app, SDL_Event* event) {
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

}

void app_render(App* app, SDL_Renderer* renderer) {
    SDL_SetRenderDrawColor(renderer, 255, 255, 255, 255);
    SDL_RenderClear(renderer);

    render_sidebar(app, renderer);
    render_nodes(app, renderer);
    render_wires(app, renderer);

    SDL_RenderPresend(renderer);
}

void app_cleanup(App* app) {

    g_list_free_full(app->nodes, free);
    g_list_free_full(app->wires, free);
    app->nodes = NULL;
    app->wires = NULL;
}

static void handle_sidebar_click(App* app, int y) {

}

static void handle_canvas_click(App* app, int x, int y) {

}