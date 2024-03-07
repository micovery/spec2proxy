package flowcallout

/**
FlowCallout:
  .name: FC-Callout
  .enabled: true
  .continueOnError: true
  DisplayName: FC-Callout
  SharedFlowBundle: FC-Callout
  Parameters:
    - Parameter:
        .name: param1
        .value: Literal
    - Parameter:
        .name: param2
        .ref: request.content
*/

type ParamT struct {
	Parameter struct {
		Name  string `json:".name" yaml:".name"`
		Value string `json:".value" yaml:".value"`
		Ref   string `json:".ref" yaml:".ref"`
	} `json:"Parameter" yaml:"Parameter"`
}

type PolicyT struct {
	Name             string   `json:".name" yaml:".name"`
	Enabled          bool     `json:".enabled" yaml:".enabled"`
	ContinueOnError  bool     `json:".continueOnError" yaml:".continueOnError"`
	DisplayName      string   `json:"DisplayName" yaml:"DisplayName"`
	SharedFlowBundle string   `json:"SharedFlowBundle" yaml:"SharedFlowBundle"`
	Parameters       []ParamT `json:"Parameters" yaml:"Parameters"`
}

type FlowCallout struct {
	FlowCallout PolicyT `policy:"true" json:"FlowCallout" yaml:"FlowCallout" `
}

func (p *FlowCallout) Name() string {
	return p.FlowCallout.Name
}

func (p *FlowCallout) Enabled() bool {
	return p.FlowCallout.Enabled
}

func (p *FlowCallout) XML() string {
	return ""
}
