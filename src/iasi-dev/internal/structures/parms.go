package structures

import "os"

type Parms struct {
	Verbose          int      // Nivel de detalle de la salida.
	All              bool     // Ejecuta también las etapas anteriores.
	Checkpoints      bool     // Usa workflows intermedios como checkpoints.
	Debug            bool     // Muestra mensajes temporales de depuración.
	Force            bool     // Fuerza la operación cuando está soportado.
	Help             bool     // Muestra la ayuda.
	Install          bool     // Instala el artefacto cuando proceda.
	Local            bool     // Hace commit sin push.
	Tolerant         bool     // Continúa cuando una operación falla.
	Message          string   // Mensaje utilizado para el commit.
	Format           string   // Formato o formatos de salida.
	Subcommand       string   // Subcomando cuando command es workflow.
	LogFile          *os.File // Handle al fichero de log de la ejecución.
	RC               *int     // Código de retorno acumulativo compartido.
	LastRC           int      // Resultado de la última operación ejecutada; lo consumen los workflows.
	Targets          []string // Raíces efectivas desde las que se descubren repositorios.
	RequestedTargets []string // Objetivos solicitados por el usuario.
	Exclusions       []string // Nombres excluidos durante el descubrimiento.
	Repos            []string // Repositorios Git descubiertos.
	BlackList        []string // Repositorios que no deben procesarse en commit.
}
