#include <gtk/gtk.h>
#include <dirent.h>
#include <glib/gstdio.h>

static GtkWidget *notes_list;
static gchar *current_path;

static void on_new_file_clicked(GtkButton *button, gpointer user_data) {
  GtkWidget *window = GTK_WIDGET(user_data);
  GtkWidget *dialog = gtk_dialog_new_with_buttons("New File",
                                                  GTK_WINDOW(window),
                                                  GTK_DIALOG_MODAL,
                                                  "_Cancel", GTK_RESPONSE_CANCEL,
                                                  "_Create", GTK_RESPONSE_ACCEPT,
                                                  NULL);
  GtkWidget *content_area = gtk_dialog_get_content_area(GTK_DIALOG(dialog));
  GtkWidget *entry = gtk_entry_new();
  gtk_container_add(GTK_CONTAINER(content_area), entry);
  gtk_widget_show_all(dialog);

  gint result = gtk_dialog_run(GTK_DIALOG(dialog));
  if (result == GTK_RESPONSE_ACCEPT) {
    const gchar *filename = gtk_entry_get_text(GTK_ENTRY(entry));
    if (filename && strlen(filename) > 0) {
      gchar *full_path = g_build_filename(current_path, filename, NULL);
      if (g_file_test(full_path, G_FILE_TEST_EXISTS)) {
        GtkWidget *error_dialog = gtk_message_dialog_new(GTK_WINDOW(window),
                                                         GTK_DIALOG_MODAL,
                                                         GTK_MESSAGE_ERROR,
                                                         GTK_BUTTONS_OK,
                                                         "File already exists!");
        gtk_dialog_run(GTK_DIALOG(error_dialog));
        gtk_widget_destroy(error_dialog);
      } else {
        g_file_set_contents(full_path, "", 0, NULL);
        populate_notes_list();
      }
      g_free(full_path);
    }
  }
  gtk_widget_destroy(dialog);
}

static void populate_notes_list() {
    GtkListStore *store = GTK_LIST_STORE(gtk_tree_view_get_model(GTK_TREE_VIEW(notes_list)));
    gtk_list_store_clear(store);

    // Add ".." entry for parent directory if not at root
    if (strcmp(current_path, "notes/") != 0) {
        GtkTreeIter iter;
        gtk_list_store_append(store, &iter);
        GtkIconTheme *icon_theme = gtk_icon_theme_get_default();
        GdkPixbuf *icon = gtk_icon_theme_load_icon(icon_theme, "go-up", 16, 0, NULL);
        gchar *parent_path = g_path_get_dirname(current_path);
        gtk_list_store_set(store, &iter, 
                           0, icon,
                           1, "..",
                           2, parent_path,
                           3, TRUE,  // Correct column and value for directory
                           -1);
        if (icon) g_object_unref(icon);
        g_free(parent_path);
    }

    // Populate the list with directory contents
    DIR *dir = opendir(current_path);
    if (dir != NULL) {
        struct dirent *entry;
        while ((entry = readdir(dir)) != NULL) {
            if (strcmp(entry->d_name, ".") == 0 || strcmp(entry->d_name, "..") == 0) {
                continue;
            }

            GtkTreeIter iter;
            gtk_list_store_append(store, &iter);
            GtkIconTheme *icon_theme = gtk_icon_theme_get_default();
            GdkPixbuf *icon;
            gchar *full_path = g_build_filename(current_path, entry->d_name, NULL);
            gboolean is_directory = (entry->d_type == DT_DIR);

            if (is_directory) {
                icon = gtk_icon_theme_load_icon(icon_theme, "folder", 16, 0, NULL);
                gtk_list_store_set(store, &iter, 
                                   0, icon,
                                   1, entry->d_name,
                                   2, full_path,
                                   3, TRUE,  // Directory
                                   -1);
            } else if (entry->d_type == DT_REG) {
                icon = gtk_icon_theme_load_icon(icon_theme, "text-x-generic", 16, 0, NULL);
                gtk_list_store_set(store, &iter, 
                                   0, icon,
                                   1, entry->d_name,
                                   2, full_path,
                                   3, FALSE,  // File, not a directory
                                   -1);
            }

            if (icon) g_object_unref(icon);
            g_free(full_path);
        }
        closedir(dir);
    }
}
// Callback when a note or directory is selected
static void on_note_selected(GtkTreeSelection *selection, gpointer data) {
    GtkWidget *text_view = GTK_WIDGET(data);  // The text editor widget
    GtkTreeModel *model;
    GtkTreeIter iter;

    if (gtk_tree_selection_get_selected(selection, &model, &iter)) {
        gchar *full_path;
        gboolean is_directory;
        gtk_tree_model_get(model, &iter, 
                           2, &full_path,    // Get full path
                           3, &is_directory, // Get directory flag
                           -1);

        if (is_directory) {
            // Navigate into the selected directory
            g_free(current_path);  // Free the old path
            current_path = g_strdup(full_path);  // Update current_path
            populate_notes_list();  // Refresh the list with new directory contents
        } else {
            // Load the selected file into the editor
            gchar *content;
            if (g_file_get_contents(full_path, &content, NULL, NULL)) {
                GtkTextBuffer *buffer = gtk_text_view_get_buffer(GTK_TEXT_VIEW(text_view));
                gtk_text_buffer_set_text(buffer, content, -1);
                g_free(content);
            }
        }
        g_free(full_path);  // Free the retrieved full path
    }
}


// Callback when the Save button is clicked
static void on_save_clicked(GtkButton *button, gpointer data) {
    GtkWidget *text_view = GTK_WIDGET(data);  // The text editor widget
    GtkTreeSelection *selection = gtk_tree_view_get_selection(GTK_TREE_VIEW(notes_list));
    GtkTreeModel *model;
    GtkTreeIter iter;

    if (gtk_tree_selection_get_selected(selection, &model, &iter)) {
        gboolean is_directory;
        gtk_tree_model_get(model, &iter, 3, &is_directory, -1);
        if (!is_directory) {  // Only save if it’s a file
            gchar *full_path;
            gtk_tree_model_get(model, &iter, 2, &full_path, -1);

            // Get the text from the editor
            GtkTextBuffer *buffer = gtk_text_view_get_buffer(GTK_TEXT_VIEW(text_view));
            GtkTextIter start, end;
            gtk_text_buffer_get_bounds(buffer, &start, &end);
            gchar *text = gtk_text_buffer_get_text(buffer, &start, &end, FALSE);

            // Save the text to the file
            g_file_set_contents(full_path, text, -1, NULL);

            g_free(text);  // Free the text
            g_free(full_path);  // Free the full path
        }
    }
}


static gboolean on_canvas_draw(GtkWidget *widget, cairo_t *cr, gpointer data) {
  cairo_set_source_rgb(cr, 0, 0, 0);
  cairo_rectangle(cr, 50, 50, 100, 100);
  cairo_stroke(cr);

  return FALSE;
}

static GtkWidget *window;
// Application activation function
static void activate(GtkApplication *app, gpointer user_data) {
    window = gtk_application_window_new(app);
    gtk_window_maximize(GTK_WINDOW(window));
    gtk_window_set_title(GTK_WINDOW(window), "Obelisk Notes Manager");
    gtk_window_set_default_size(GTK_WINDOW(window), 800, 600);

    GtkWidget *notebook = gtk_notebook_new();
    gtk_notebook_set_tab_pos(GTK_NOTEBOOK(notebook), GTK_POS_LEFT);

    // Notes Tab
    GtkWidget *notes_paned = gtk_paned_new(GTK_ORIENTATION_HORIZONTAL);

    // Left side: Tree view and New File button
    GtkWidget *left_box = gtk_box_new(GTK_ORIENTATION_VERTICAL, 0);
    notes_list = gtk_tree_view_new();
    GtkListStore *store = gtk_list_store_new(4, GDK_TYPE_PIXBUF, G_TYPE_STRING, G_TYPE_STRING, G_TYPE_BOOLEAN);
    gtk_tree_view_set_model(GTK_TREE_VIEW(notes_list), GTK_TREE_MODEL(store));
    g_object_unref(store);

    GtkTreeViewColumn *column = gtk_tree_view_column_new();
    gtk_tree_view_column_set_title(column, "Notes");

    GtkCellRenderer *pixbuf_renderer = gtk_cell_renderer_pixbuf_new();
    gtk_tree_view_column_pack_start(column, pixbuf_renderer, FALSE);
    gtk_tree_view_column_add_attribute(column, pixbuf_renderer, "pixbuf", 0);

    GtkCellRenderer *text_renderer = gtk_cell_renderer_text_new();
    gtk_tree_view_column_pack_start(column, text_renderer, TRUE);
    gtk_tree_view_column_add_attribute(column, text_renderer, "text", 1);

    gtk_tree_view_append_column(GTK_TREE_VIEW(notes_list), column);

    // Set fixed width for the column
    gtk_tree_view_column_set_sizing(column, GTK_TREE_VIEW_COLUMN_FIXED);
    gtk_tree_view_column_set_fixed_width(column, 200);

    gtk_box_pack_start(GTK_BOX(left_box), notes_list, TRUE, TRUE, 0);

    GtkWidget *new_file_button = gtk_button_new_with_label("New File");
    gtk_box_pack_start(GTK_BOX(left_box), new_file_button, FALSE, FALSE, 0);
    g_signal_connect(new_file_button, "clicked", G_CALLBACK(on_new_file_clicked), window);

    gtk_paned_pack1(GTK_PANED(notes_paned), left_box, TRUE, FALSE);

    // Right side: Editor (your existing code goes here)
    // gtk_paned_pack2(GTK_PANED(notes_paned), editor_box, TRUE, FALSE);

    // Populate the list initially
    populate_notes_list();

    GtkWidget *notes_label = gtk_label_new("Notes");
    gtk_notebook_append_page(GTK_NOTEBOOK(notebook), notes_paned, notes_label);

    // Rest of your code (Canvas Tab, Settings Tab, etc.)

    gtk_container_add(GTK_CONTAINER(window), notebook);
    gtk_widget_show_all(window);
}


int main(int argc, char **argv) {
    // Use G_APPLICATION_FLAGS_NONE instead of G_APPLICATION_DEFAULT_FLAGS
    GtkApplication *app = gtk_application_new("Obelisk.Notes.Manager", G_APPLICATION_FLAGS_NONE);
  current_path = g_strdup("notes/");
    g_signal_connect(app, "activate", G_CALLBACK(activate), NULL);
    int status = g_application_run(G_APPLICATION(app), argc, argv);
    g_object_unref(app);
    return status;
}
