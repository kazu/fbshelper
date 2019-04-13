package example

type MoveStateMessage struct {
    
        Gps Gps
    
        ObjectId int64
    
        Velocity Vec3
    
    Dirty bool
    cache *example.MoveState
}

func NewMoveStateMessageFromFbs(fbs *example.MoveState) (m *MoveStateMessage) {
    return &MoveStateMessage{
         
            Gps: fbs.Gps(),
         
            ObjectId: fbs.ObjectId(),
         
            Velocity: fbs.Velocity(),
         
         Dirty: false,
         cache: fbs,   
    }
}

func (m *MoveStateMessage) ApplyUpdated() bool {
    
        if m.Gps != m.cache.Gps() {
             return m.cache.MutateGps(m.Gps)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    
        if m.ObjectId != m.cache.ObjectId() {
             return m.cache.MutateObjectId(m.ObjectId)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    
        if m.Velocity != m.cache.Velocity() {
             return m.cache.MutateVelocity(m.Velocity)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    

    return true
}

func (m *MoveStateMessage) ToFbs(check bool) (fbs *example.MoveState) {
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
