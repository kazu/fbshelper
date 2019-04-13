package example

type CollisionMessage struct {
    
        MoveStates []MoveState
    
    Dirty bool
    cache *example.Collision
}

func NewCollisionMessageFromFbs(fbs *example.Collision) (m *CollisionMessage) {
    return &CollisionMessage{
         
            MoveStates: fbs.MoveStates(),
         
         Dirty: false,
         cache: fbs,   
    }
}

func (m *CollisionMessage) ApplyUpdated() bool {
    
        if m.MoveStates != m.cache.MoveStates() {
             return m.cache.MutateMoveStates(m.MoveStates)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    

    return true
}

func (m *CollisionMessage) ToFbs(check bool) (fbs *example.Collision) {
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
