package example

type TimeMessage struct {
    
        Nano int64
    
        Unix int64
    
    Dirty bool
    cache *example.Time
}

func NewTimeMessageFromFbs(fbs *example.Time) (m *TimeMessage) {
    return &TimeMessage{
         
            Nano: fbs.Nano(),
         
            Unix: fbs.Unix(),
         
         Dirty: false,
         cache: fbs,   
    }
}

func (m *TimeMessage) ApplyUpdated() bool {
    
        if m.Nano != m.cache.Nano() {
             return m.cache.MutateNano(m.Nano)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    
        if m.Unix != m.cache.Unix() {
             return m.cache.MutateUnix(m.Unix)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    

    return true
}

func (m *TimeMessage) ToFbs(check bool) (fbs *example.Time) {
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
