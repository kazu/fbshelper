package example

type RootMessage struct {
    
        Action Action
    
        ActionTime Time
    
        MoveState MoveState
    
        Time Time
    
    Dirty bool
    cache *example.Root
}

func NewRootMessageFromFbs(fbs *example.Root) (m *RootMessage) {
    return &RootMessage{
         
            Action: fbs.Action(),
         
            ActionTime: fbs.ActionTime(),
         
            MoveState: fbs.MoveState(),
         
            Time: fbs.Time(),
         
         Dirty: false,
         cache: fbs,   
    }
}

func (m *RootMessage) ApplyUpdated() bool {
    
        if m.Action != m.cache.Action() {
             return m.cache.MutateAction(m.Action)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    
        if m.ActionTime != m.cache.ActionTime() {
             return m.cache.MutateActionTime(m.ActionTime)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    
        if m.MoveState != m.cache.MoveState() {
             return m.cache.MutateMoveState(m.MoveState)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    
        if m.Time != m.cache.Time() {
             return m.cache.MutateTime(m.Time)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    

    return true
}

func (m *RootMessage) ToFbs(check bool) (fbs *example.Root) {
    if !check && m.cache != nil {
        return m.cache
    }
    if m.ApplyUpdated() {
        return m.cache
    }
    m.Dirty = true
    // 作成処理
    /*
          fbs = ...   
          m.cache = fbs
    */

    return m.cache
}
