package server

var (
	fastEnd       = "/operation/fast-end"
	slowEnd       = "/operation/slow-end"
	singleCall    = "/operation/single-call"
	multipleCalls = "/operation/multiple-calls"

	routePaths = []string{
		fastEnd,
		slowEnd,
		singleCall,
		multipleCalls,
	}
)

func (s *server) routes() {
	s.router.HandleFunc(fastEnd,
		s.middlewareLogger(
			s.middlewareTrace(
				s.fastEnd())))

	s.router.HandleFunc(slowEnd,
		s.middlewareLogger(
			s.middlewareTrace(
				s.slowEnd())))

	s.router.HandleFunc(singleCall,
		s.middlewareLogger(
			s.middlewareTrace(
				s.singleCall())))

	s.router.HandleFunc(multipleCalls,
		s.middlewareLogger(
			s.middlewareTrace(
				s.multipleCalls())))
}
