import React from "react";
import DashboardMain from '../../src/components/DashboardMain/DashboardMain';
import {useRouter} from "next/router";

export default function Index(props) {

    const router = useRouter();
    React.useEffect(() => {
        router.push('/dashboard/application');
    }, [])
    return (
        <DashboardMain>
        </DashboardMain>
    );
}