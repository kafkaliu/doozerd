package paxos

const (
    lRnd = iota
    lValue
    lNumParts
)

func learner(quorum int, ins chan Msg) string {
    var round uint64 = 0
    votes := make(map[string]int) // maps values to number of votes
    voted := make(map[uint64]bool) // maps values to number of votes

    update := func(in Msg) string {
        defer swallowContinue()

        parts := splitExactly(in.body, lNumParts) // e.g. 1:xxx

        if in.cmd != "VOTE" {
            return ""
        }

        mRound := dtoui64(parts[lRnd])

        v := parts[lValue]

        switch {
        case mRound < round:
            return ""
        case mRound > round:
            round = mRound
            votes = make(map[string]int)
            voted = make(map[uint64]bool)
            fallthrough
        case mRound == round:
            if voted[in.from] {
                return ""
            }
            votes[v]++
            voted[in.from] = true

            if votes[v] >= quorum {
                return v // winner!
            }
        }

        return ""
    }

    for in := range ins {
        v := update(in)
        if v != "" {
            return v
        }
    }
    panic("can't happen")
}
