import * as React from 'react';
import {Header} from '../src/components/Header/Header';
import {Hero} from '../src/components/Hero/Hero';
import {Features} from '../src/components/Features/Features';
import {TryIt} from '../src/components/TryIt/TryIt';
import { Footer } from '../src/components/Footer/Footer';

export default function Index() {
  return (
    <div>
        <div style={{height: '100vh'}}>
            <Header />
            <Hero />
        </div>
        <Features />
        <TryIt />
        <Footer />
    </div>
  );
}
