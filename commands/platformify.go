package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/models"
	"github.com/platformsh/platformify/internal/question"
	"github.com/platformsh/platformify/internal/questionnaire"
	"github.com/platformsh/platformify/platformifiers"
)

// PlatformifyCmd represents the base Platformify command when called without any subcommands
var PlatformifyCmd = &cobra.Command{
	Use:   "platformify",
	Short: "Platfomrify your application, and deploy it to the Platform.sh",
	Long: `Platformify your application, creating all the needed files
for it to be deployed to Platform.sh.

This will create the needed YAML files for both your application and your
services, choosing from a variety of stacks or simple runtimes.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		answers := models.NewAnswers()
		ctx := models.ToContext(cmd.Context(), answers)
		ctx = colors.ToContext(
			ctx,
			cmd.OutOrStderr(),
			cmd.ErrOrStderr(),
		)
		q := questionnaire.New(
			&question.WorkingDirectory{},
			&question.Welcome{},
			&question.Stack{},
			&question.Type{},
			&question.DependencyManager{},
			&question.HalfWay{},
			&question.Name{},
			&question.ApplicationRoot{},
			&question.Environment{},
			&question.BuildSteps{},
			&question.DeployCommand{},
			&question.ListenInterface{},
			&question.WebCommand{},
			&question.AlmostDone{},
			&question.Services{},
			&question.Done{},
		)
		err := q.AskQuestions(ctx)
		if err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), colors.Colorize(colors.ErrorCode, err.Error()))
			return err
		}

		pfier, err := platformifiers.NewPlatformifier(answers)
		if err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), colors.Colorize(colors.ErrorCode, err.Error()))
			return fmt.Errorf("creating platformifier failed: %s", err)
		}

		if err := pfier.Platformify(ctx); err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), colors.Colorize(colors.ErrorCode, err.Error()))
			return fmt.Errorf("could not platformify project: %w", err)
		}

		return nil
	},
}