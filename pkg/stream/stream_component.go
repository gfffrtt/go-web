package stream

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

type StreamComponent struct {
	ID        string
	Component templ.Component
}

func NewStreamComponent(id string, component templ.Component) *StreamComponent {
	return &StreamComponent{
		ID:        id,
		Component: component,
	}
}

func (sc *StreamComponent) Render(ctx context.Context, w http.ResponseWriter) error {
	err := sc.Component.Render(ctx, w)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(fmt.Sprintf(`
		<script>
			(() => {
				const content = document.querySelector('template[data-content="%s"]');
				const loading = document.querySelector('[data-loading="%s"]');
				const error = document.querySelector('[data-error="%s"]');
				if (!content) {
					loading.replaceWith(error.content);
				}
				loading.replaceWith(content.content);
			})();
		</script>
	`, sc.ID, sc.ID, sc.ID)))
	if err != nil {
		return err
	}
	return nil
}
