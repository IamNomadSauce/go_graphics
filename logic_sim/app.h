fndef APP_H
#define APP_H

#include "types.h"

void app_init(App*app);
void app_handle_event(App* app, SDL_Event* event);
void app_update(App* app);
void app_render(App* app, SDL_Renderer* renderer);
void app_cleanup(App* app);

#endif