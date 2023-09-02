use super::{Sys, SystemArchitecture};

impl SystemArchitecture for Sys {
    fn get_cpu_architecture() -> Option<String> {
        // https://gitlab.com/x86-psABIs/x86-64-ABI/-/jobs/artifacts/master/raw/x86-64-ABI/abi.pdf?job=build
        // https://gitlab.com/x86-psABIs/x86-64-ABI/-/blob/master/x86-64-ABI/low-level-sys-info.tex
        // https://unix.stackexchange.com/questions/43539/what-do-the-flags-in-proc-cpuinfo-mean/43540#43540

        let cpuid = raw_cpuid::CpuId::new();

        let finfo = cpuid.get_feature_info()?;
        let efinfo = cpuid.get_extended_feature_info()?;
        let epfi = cpuid.get_extended_processor_and_feature_identifiers()?;

        let is_i686 = finfo.has_cmov();

        let is_x86_64 = is_i686
            && epfi.has_64bit_mode()
            && finfo.has_cmpxchg8b()
            && finfo.has_fpu()
            && finfo.has_fxsave_fxstor()
            && finfo.has_mmx()
            && epfi.has_syscall_sysret()
            && finfo.has_sse()
            && finfo.has_sse2();

        let is_x86_64_v2 = is_x86_64
            && finfo.has_cmpxchg16b()
            && epfi.has_lahf_sahf()
            && finfo.has_popcnt()
            && finfo.has_sse3()
            && finfo.has_sse41()
            && finfo.has_sse42()
            && finfo.has_ssse3();

        let is_x86_64_v3 = is_x86_64_v2
            && finfo.has_avx()
            && efinfo.has_avx2()
            && efinfo.has_bmi1()
            && efinfo.has_bmi2()
            && finfo.has_f16c()
            && finfo.has_fma()
            && epfi.has_lzcnt()
            && finfo.has_movbe()
            && finfo.has_oxsave();

        let is_x86_64_v4 = is_x86_64_v3
            && efinfo.has_avx512f()
            && efinfo.has_avx512bw()
            && efinfo.has_avx512cd()
            && efinfo.has_avx512dq()
            && efinfo.has_avx512vl();

        if is_x86_64_v4 {
            return Some("x86_64_v4".to_string());
        } else if is_x86_64_v3 {
            return Some("x86_64_v3".to_string());
        } else if is_x86_64_v2 {
            return Some("x86_64_v2".to_string());
        } else if is_x86_64 {
            return Some("x86_64".to_string());
        } else if is_i686 {
            return Some("i686".to_string());
        }

        None
    }
}
