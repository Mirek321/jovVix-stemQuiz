package templates

import "fmt"

// permissionLabels maps permission enum values to user-facing Czech labels.
var permissionLabels = map[string]string{
	"read":  "čtení",
	"write": "úpravy",
	"share": "sdílení",
}

// GenerateQuizShareEmail generates the HTML content for the quiz share email
func GenerateQuizShareEmail(quizLink, sharedBy string, permissions string) string {
	permissionLabel := permissions
	if label, ok := permissionLabels[permissions]; ok {
		permissionLabel = label
	}
	permissionsText := fmt.Sprintf("Udělená oprávnění: <strong>%s</strong>", permissionLabel)

	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="cs">
		<head>
			<style>
				body { font-family: Arial, sans-serif; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #eaeaea; border-radius: 8px; }
				.button { display: inline-block; padding: 10px 20px; margin: 20px 0; background-color: #007bff; color: #ffffff; text-decoration: none; border-radius: 5px; }
				.footer { font-size: 12px; color: #888; margin-top: 20px; }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Ahoj,</h2>
				<p><strong>%s</strong> tě zve ke kvízu.</p>
				<p>%s</p>
				<p>Do kvízu se dostaneš kliknutím na odkaz níže:</p>
				<p><a href="%s" class="button" style="color:#ffffff !important; text-decoration:none;">Otevřít kvíz</a></p>
				<p>Pokud tlačítko nefunguje, zkopíruj tento odkaz a vlož ho do prohlížeče:</p>
				<p><a href="%s">%s</a></p>
				<p class="footer">Tuto pozvánku ti poslal(a) %s prostřednictvím kvízové platformy.</p>
			</div>
		</body>
		</html>`, sharedBy, permissionsText, quizLink, quizLink, quizLink, sharedBy)
}
