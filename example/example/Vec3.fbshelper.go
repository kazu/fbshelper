package example

type Vec3Message struct {
    
        X float64
    
        Y float64
    
        Z float64
    
    Dirty bool
    cache *example.Vec3
}

func NewVec3MessageFromFbs(fbs *example.Vec3) (m *Vec3Message) {
    return &Vec3Message{
         
            X: fbs.X(),
         
            Y: fbs.Y(),
         
            Z: fbs.Z(),
         
         Dirty: false,
         cache: fbs,   
    }
}

func (m *Vec3Message) ApplyUpdated() bool {
    
        if m.X != m.cache.X() {
             return m.cache.MutateX(m.X)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    
        if m.Y != m.cache.Y() {
             return m.cache.MutateY(m.Y)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    
        if m.Z != m.cache.Z() {
             return m.cache.MutateZ(m.Z)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    

    return true
}

func (m *Vec3Message) ToFbs(check bool) (fbs *example.Vec3) {
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
